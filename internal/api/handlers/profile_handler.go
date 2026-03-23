package handlers

import (
	"net/http"

	"enterprise-orbit/internal/models"
	"enterprise-orbit/pkg/auth"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ProfileHandler struct {
	db *gorm.DB
}

func NewProfileHandler(db *gorm.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	response := models.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Nickname:   user.Nickname,
		Phone:      user.Phone,
		Avatar:     user.Avatar,
		Department: user.Department,
		Position:   user.Position,
		Role:       user.Role,
		Status:     user.Status,
		LastLogin:  user.LastLogin,
		CreatedAt:  user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	updates := map[string]interface{}{
		"nickname":   req.Nickname,
		"phone":      req.Phone,
		"avatar":     req.Avatar,
		"department": req.Department,
		"position":   req.Position,
	}

	if err := h.db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	h.db.First(&user, userID)
	response := models.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Nickname:   user.Nickname,
		Phone:      user.Phone,
		Avatar:     user.Avatar,
		Department: user.Department,
		Position:   user.Position,
		Role:       user.Role,
		Status:     user.Status,
		LastLogin:  user.LastLogin,
		CreatedAt:  user.CreatedAt,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProfileHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := h.db.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	filename := auth.GenerateAvatarFilename(userID.(uint), file.Filename)
	savePath := "./uploads/avatars/" + filename

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	avatarURL := "/uploads/avatars/" + filename
	h.db.Model(&models.User{}).Where("id = ?", userID).Update("avatar", avatarURL)

	c.JSON(http.StatusOK, gin.H{"url": avatarURL})
}
