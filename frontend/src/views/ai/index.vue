<template>
  <div class="ai-assistant">
    <div class="page-header">
      <div class="page-title">
        <h1>{{ $t('ai.title') }}</h1>
        <p class="subtitle">{{ $t('ai.subtitle') }}</p>
      </div>
    </div>

    <div class="chat-container glass">
      <div class="chat-messages" ref="messagesContainer">
        <div v-if="messages.length === 0" class="welcome-screen">
          <div class="welcome-icon">
            <svg viewBox="0 0 80 80" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="40" cy="40" r="36" stroke="url(#grad1)" stroke-width="2"/>
              <circle cx="40" cy="40" r="16" fill="url(#grad1)"/>
              <circle cx="40" cy="16" r="6" fill="#06b6d4"/>
              <circle cx="64" cy="40" r="6" fill="#8b5cf6"/>
              <circle cx="40" cy="64" r="6" fill="#06b6d4"/>
              <circle cx="16" cy="40" r="6" fill="#8b5cf6"/>
              <defs>
                <linearGradient id="grad1" x1="0" y1="0" x2="80" y2="80">
                  <stop offset="0%" stop-color="#06b6d4"/>
                  <stop offset="100%" stop-color="#8b5cf6"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <h2>{{ $t('ai.welcomeTitle') }}</h2>
          <p>{{ $t('ai.welcomeDesc') }}</p>
          <div class="quick-actions">
            <button v-for="action in quickActions" :key="action" class="quick-action-btn" @click="sendQuickAction(action)">
              {{ action }}
            </button>
          </div>
        </div>

        <div v-for="(message, index) in messages" :key="index" class="message" :class="message.role">
          <div class="message-avatar">
            <div v-if="message.role === 'user'" class="avatar user-avatar">
              <el-icon><User /></el-icon>
            </div>
            <div v-else class="avatar ai-avatar">
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5"/>
                <circle cx="12" cy="12" r="4" fill="currentColor"/>
                <circle cx="12" cy="4" r="2" fill="currentColor"/>
                <circle cx="20" cy="12" r="2" fill="currentColor"/>
                <circle cx="12" cy="20" r="2" fill="currentColor"/>
                <circle cx="4" cy="12" r="2" fill="currentColor"/>
              </svg>
            </div>
          </div>
          <div class="message-content">
            <div class="message-header">
              <span class="message-role">{{ message.role === 'user' ? $t('ai.you') : $t('ai.assistant') }}</span>
              <span class="message-time">{{ formatTime(message.timestamp) }}</span>
            </div>
            <div class="message-text" v-html="renderMarkdown(message.content)"></div>
          </div>
        </div>

        <div v-if="isLoading" class="message assistant">
          <div class="message-avatar">
            <div class="avatar ai-avatar thinking">
              <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="1.5"/>
                <circle cx="12" cy="12" r="4" fill="currentColor"/>
              </svg>
            </div>
          </div>
          <div class="message-content">
            <div class="thinking-indicator">
              <span></span>
              <span></span>
              <span></span>
            </div>
          </div>
        </div>
      </div>

      <div class="chat-input-area">
        <div class="input-container">
          <el-input
            v-model="inputMessage"
            type="textarea"
            :rows="1"
            :autosize="{ minRows: 1, maxRows: 4 }"
            :placeholder="$t('ai.inputPlaceholder')"
            @keydown.enter.exact.prevent="sendMessage"
            :disabled="isLoading"
            class="chat-input"
          />
          <el-button
            type="primary"
            :loading="isLoading"
            :disabled="!inputMessage.trim()"
            @click="sendMessage"
            class="send-btn"
          >
            <el-icon v-if="!isLoading"><Promotion /></el-icon>
          </el-button>
        </div>
        <div class="input-hint">
          <span>{{ $t('ai.inputHint') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { User, Promotion } from '@element-plus/icons-vue'
import axios from 'axios'
import { marked } from 'marked'

const { t } = useI18n()

interface Message {
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
}

const messages = ref<Message[]>([])
const inputMessage = ref('')
const isLoading = ref(false)
const messagesContainer = ref<HTMLElement | null>(null)

const quickActions = [
  t('ai.quickAction1'),
  t('ai.quickAction2'),
  t('ai.quickAction3'),
  t('ai.quickAction4')
]

const formatTime = (date: Date) => {
  return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

const renderMarkdown = (content: string) => {
  return marked(content)
}

const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
  })
}

const sendMessage = async () => {
  if (!inputMessage.value.trim() || isLoading.value) return

  const userMessage = inputMessage.value.trim()
  inputMessage.value = ''

  messages.value.push({
    role: 'user',
    content: userMessage,
    timestamp: new Date()
  })
  scrollToBottom()

  isLoading.value = true

  try {
    const token = localStorage.getItem('token')
    const chatMessages = messages.value.map(m => ({
      role: m.role,
      content: m.content
    }))

    const response = await fetch('/api/ai/chat/stream', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ messages: chatMessages })
    })

    if (!response.ok) {
      throw new Error('Failed to connect to AI service')
    }

    const reader = response.body?.getReader()
    if (!reader) {
      throw new Error('No reader available')
    }

    const assistantMessage: Message = {
      role: 'assistant',
      content: '',
      timestamp: new Date()
    }
    messages.value.push(assistantMessage)
    scrollToBottom()

    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      assistantMessage.content = buffer
      scrollToBottom()
    }

  } catch (error) {
    console.error('Chat error:', error)
    messages.value.push({
      role: 'assistant',
      content: t('ai.errorMessage'),
      timestamp: new Date()
    })
  } finally {
    isLoading.value = false
    scrollToBottom()
  }
}

const sendQuickAction = (action: string) => {
  inputMessage.value = action
  sendMessage()
}

onMounted(() => {
  scrollToBottom()
})
</script>

<style scoped>
.ai-assistant {
  max-width: 1000px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
}

.page-title h1 {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.subtitle {
  color: var(--color-text-secondary);
  font-size: 14px;
  margin-top: 8px;
}

.chat-container {
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border);
  background: rgba(10, 15, 28, 0.6);
  backdrop-filter: blur(20px);
  display: flex;
  flex-direction: column;
  height: calc(100vh - 200px);
  min-height: 500px;
  position: relative;
  overflow: hidden;
}

.chat-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, var(--color-accent), var(--color-accent-secondary));
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  scroll-behavior: smooth;
}

.welcome-screen {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  text-align: center;
  padding: 40px;
}

.welcome-icon {
  width: 80px;
  height: 80px;
  margin-bottom: 24px;
  animation: float 3s ease-in-out infinite;
}

.welcome-icon svg {
  width: 100%;
  height: 100%;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.welcome-screen h2 {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 12px 0;
}

.welcome-screen p {
  color: var(--color-text-secondary);
  font-size: 14px;
  margin: 0 0 32px 0;
  max-width: 400px;
}

.quick-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  justify-content: center;
  max-width: 500px;
}

.quick-action-btn {
  padding: 10px 16px;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: rgba(255, 255, 255, 0.02);
  color: var(--color-text-secondary);
  font-size: 13px;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.quick-action-btn:hover {
  border-color: var(--color-accent);
  color: var(--color-accent);
  background: rgba(6, 182, 212, 0.1);
}

.message {
  display: flex;
  gap: 16px;
  margin-bottom: 24px;
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.message.user {
  flex-direction: row-reverse;
}

.message-avatar {
  flex-shrink: 0;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-avatar {
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-secondary) 100%);
  color: white;
}

.ai-avatar {
  background: rgba(139, 92, 246, 0.2);
  border: 1px solid rgba(139, 92, 246, 0.3);
  color: var(--color-accent-secondary);
}

.ai-avatar svg {
  width: 20px;
  height: 20px;
}

.ai-avatar.thinking svg {
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.message-content {
  flex: 1;
  max-width: 80%;
}

.message.user .message-content {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.message.user .message-header {
  flex-direction: row-reverse;
}

.message-role {
  font-size: 13px;
  font-weight: 500;
  color: var(--color-text-primary);
}

.message-time {
  font-size: 11px;
  color: var(--color-text-muted);
}

.message-text {
  padding: 16px 20px;
  border-radius: var(--radius-md);
  font-size: 14px;
  line-height: 1.7;
  word-break: break-word;
}

.message.user .message-text {
  background: linear-gradient(135deg, rgba(6, 182, 212, 0.15) 0%, rgba(139, 92, 246, 0.15) 100%);
  border: 1px solid rgba(6, 182, 212, 0.2);
  color: var(--color-text-primary);
}

.message.assistant .message-text {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--color-border);
  color: var(--color-text-primary);
}

.message-text :deep(p) {
  margin: 0 0 12px 0;
}

.message-text :deep(p:last-child) {
  margin-bottom: 0;
}

.message-text :deep(code) {
  background: rgba(0, 0, 0, 0.3);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Fira Code', monospace;
  font-size: 13px;
}

.message-text :deep(pre) {
  background: rgba(0, 0, 0, 0.4);
  padding: 16px;
  border-radius: var(--radius-sm);
  overflow-x: auto;
  margin: 12px 0;
}

.message-text :deep(pre code) {
  background: none;
  padding: 0;
}

.message-text :deep(ul), .message-text :deep(ol) {
  margin: 8px 0;
  padding-left: 24px;
}

.message-text :deep(li) {
  margin: 4px 0;
}

.message-text :deep(strong) {
  color: var(--color-accent);
}

.thinking-indicator {
  display: flex;
  gap: 4px;
  padding: 16px 20px;
}

.thinking-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--color-accent);
  animation: bounce 1.4s ease-in-out infinite;
}

.thinking-indicator span:nth-child(1) { animation-delay: 0s; }
.thinking-indicator span:nth-child(2) { animation-delay: 0.2s; }
.thinking-indicator span:nth-child(3) { animation-delay: 0.4s; }

@keyframes bounce {
  0%, 80%, 100% { transform: translateY(0); }
  40% { transform: translateY(-8px); }
}

.chat-input-area {
  padding: 20px 24px;
  border-top: 1px solid var(--color-border);
  background: rgba(0, 0, 0, 0.2);
}

.input-container {
  display: flex;
  gap: 12px;
  align-items: flex-end;
}

.chat-input {
  flex: 1;
}

.chat-input :deep(.el-textarea__inner) {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-primary);
  font-size: 14px;
  padding: 12px 16px;
  resize: none;
  transition: all var(--transition-fast);
}

.chat-input :deep(.el-textarea__inner:focus) {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(6, 182, 212, 0.1);
}

.chat-input :deep(.el-textarea__inner::placeholder) {
  color: var(--color-text-muted);
}

.send-btn {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, var(--color-accent) 0%, var(--color-accent-secondary) 100%);
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
}

.send-btn:hover:not(:disabled) {
  transform: scale(1.05);
  box-shadow: 0 4px 20px rgba(6, 182, 212, 0.4);
}

.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.send-btn .el-icon {
  font-size: 20px;
  color: white;
}

.input-hint {
  margin-top: 8px;
  font-size: 12px;
  color: var(--color-text-muted);
  text-align: center;
}
</style>
