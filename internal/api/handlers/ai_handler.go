package handlers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	apiKey string
}

func NewAIHandler(apiKey string) *AIHandler {
	return &AIHandler{apiKey: apiKey}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model      string        `json:"model"`
	Messages   []ChatMessage `json:"messages"`
	Stream     bool          `json:"stream"`
	MaxTokens  int           `json:"max_tokens,omitempty"`
	Temperature float64      `json:"temperature,omitempty"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type StreamChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

type ChatHistoryRequest struct {
	Messages []ChatMessage `json:"messages" binding:"required"`
}

func (h *AIHandler) ChatStream(c *gin.Context) {
	var req ChatHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	systemPrompt := ChatMessage{
		Role:    "system",
		Content: "你是Enterprise Orbit企业关系管理系统的AI助手。你精通客户关系管理、合同管理、销售策略和商业分析。请用专业、友好的语气回答用户问题，帮助用户更好地管理企业关系。回答时请使用markdown格式，使内容更加清晰易读。",
	}

	messages := make([]ChatMessage, 0, len(req.Messages)+1)
	messages = append(messages, systemPrompt)
	messages = append(messages, req.Messages...)

	chatReq := ChatRequest{
		Model:       "qwen-plus",
		Messages:    messages,
		Stream:      true,
		MaxTokens:   2048,
		Temperature: 0.7,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
		return
	}

	httpReq, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+h.apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to AI service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{"error": string(body)})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	reader := bufio.NewReader(resp.Body)
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Streaming not supported"})
		return
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !strings.HasPrefix(line, "data:") {
			continue
		}

		data := strings.TrimPrefix(line, "data:")
		data = strings.TrimSpace(data)

		if data == "[DONE]" {
			c.SSEvent("message", "[DONE]")
			flusher.Flush()
			break
		}

		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				fmt.Fprintf(c.Writer, "data: %s\n\n", content)
				flusher.Flush()
			}
		}
	}
}

func (h *AIHandler) Chat(c *gin.Context) {
	var req ChatHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	systemPrompt := ChatMessage{
		Role:    "system",
		Content: "你是Enterprise Orbit企业关系管理系统的AI助手。你精通客户关系管理、合同管理、销售策略和商业分析。请用专业、友好的语气回答用户问题，帮助用户更好地管理企业关系。回答时请使用markdown格式，使内容更加清晰易读。",
	}

	messages := make([]ChatMessage, 0, len(req.Messages)+1)
	messages = append(messages, systemPrompt)
	messages = append(messages, req.Messages...)

	chatReq := ChatRequest{
		Model:       "qwen-plus",
		Messages:    messages,
		Stream:      false,
		MaxTokens:   2048,
		Temperature: 0.7,
	}

	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
		return
	}

	httpReq, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+h.apiKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to AI service"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": string(body)})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}
