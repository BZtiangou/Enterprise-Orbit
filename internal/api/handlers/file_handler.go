package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"enterprise-orbit/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FileHandler struct {
	db *gorm.DB
}

func NewFileHandler(db *gorm.DB) *FileHandler {
	return &FileHandler{db: db}
}

func (h *FileHandler) UploadContractFile(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var contract models.Contract
	if err := h.db.First(&contract, contractID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	category := c.DefaultPostForm("category", "contract")
	description := c.PostForm("description")
	isShared := c.PostForm("is_shared") == "true"

	allowedTypes := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".png":  true,
		".jpg":  true,
		".jpeg": true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedTypes[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
		return
	}

	if file.Size > 50*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 50MB limit"})
		return
	}

	uploadDir := fmt.Sprintf("./uploads/contracts/%d", contractID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, file.Filename)
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		userID = c.GetUint("userID")
	}

	contractFile := models.ContractFile{
		ContractID:  uint(contractID),
		FileName:    file.Filename,
		FilePath:    fmt.Sprintf("/uploads/contracts/%d/%s", contractID, filename),
		FileSize:    file.Size,
		FileType:    ext,
		Category:    category,
		Description: description,
		UploadedBy:  userID,
		IsShared:    isShared,
		CreatedAt:   time.Now(),
	}

	if err := h.db.Create(&contractFile).Error; err != nil {
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": contractFile})
}

func (h *FileHandler) GetContractFiles(c *gin.Context) {
	contractID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var files []models.ContractFile
	if err := h.db.Where("contract_id = ?", contractID).Order("created_at DESC").Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": files})
}

func (h *FileHandler) DeleteContractFile(c *gin.Context) {
	fileID, err := strconv.ParseUint(c.Param("file_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	var file models.ContractFile
	if err := h.db.First(&file, fileID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	if err := os.Remove("." + file.FilePath); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	if err := h.db.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func (h *FileHandler) UpdateFileShareStatus(c *gin.Context) {
	fileID, err := strconv.ParseUint(c.Param("file_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	var req struct {
		IsShared bool `json:"is_shared"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.db.Model(&models.ContractFile{}).Where("id = ?", fileID).Update("is_shared", req.IsShared).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File updated successfully"})
}
