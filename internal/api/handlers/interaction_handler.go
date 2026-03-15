package handlers

import (
	"net/http"
	"strconv"
	"time"

	"enterprise-orbit/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InteractionHandler struct {
	db *gorm.DB
}

func NewInteractionHandler(db *gorm.DB) *InteractionHandler {
	return &InteractionHandler{db: db}
}

type CreateInteractionRequest struct {
	CustomerID      uint    `json:"customer_id" binding:"required"`
	ContactID       *uint   `json:"contact_id"`
	InteractionType string  `json:"interaction_type" binding:"required"`
	InteractionDate string  `json:"interaction_date" binding:"required"`
	Topic           string  `json:"topic"`
	ContentSummary  string  `json:"content_summary"`
	Outcome         string  `json:"outcome"`
	NextAction      string  `json:"next_action"`
	NextActionDate  *string `json:"next_action_date"`
	ImportanceLevel int     `json:"importance_level"`
}

type UpdateInteractionRequest struct {
	ContactID       *uint   `json:"contact_id"`
	InteractionType string  `json:"interaction_type"`
	InteractionDate string  `json:"interaction_date"`
	Topic           string  `json:"topic"`
	ContentSummary  string  `json:"content_summary"`
	Outcome         string  `json:"outcome"`
	NextAction      string  `json:"next_action"`
	NextActionDate  *string `json:"next_action_date"`
	ImportanceLevel int     `json:"importance_level"`
}

func (h *InteractionHandler) GetInteractions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	interactionType := c.Query("type")
	customerID := c.Query("customer_id")

	var interactions []models.InteractionLog
	var total int64

	query := h.db.Model(&models.InteractionLog{}).Preload("Customer").Preload("Tags")

	if search != "" {
		query = query.Where("topic LIKE ? OR content_summary LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if interactionType != "" {
		query = query.Where("interaction_type = ?", interactionType)
	}
	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("interaction_date DESC").Find(&interactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch interactions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": interactions,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func (h *InteractionHandler) GetInteraction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interaction ID"})
		return
	}

	var interaction models.InteractionLog
	if err := h.db.Preload("Customer").
		Preload("Contact").
		Preload("Attachments").
		Preload("Tags").
		First(&interaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interaction})
}

func (h *InteractionHandler) CreateInteraction(c *gin.Context) {
	var req CreateInteractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.GetUint("userID")

	interactionDate, _ := time.Parse("2006-01-02", req.InteractionDate)

	interaction := models.InteractionLog{
		CustomerID:      req.CustomerID,
		ContactID:       req.ContactID,
		InteractionType: req.InteractionType,
		InteractionDate: interactionDate,
		Topic:           req.Topic,
		ContentSummary:  req.ContentSummary,
		Outcome:         req.Outcome,
		NextAction:      req.NextAction,
		ImportanceLevel: req.ImportanceLevel,
		CreatedBy:       userID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if req.NextActionDate != nil {
		nextDate, _ := time.Parse("2006-01-02", *req.NextActionDate)
		interaction.NextActionDate = &nextDate
	}

	if err := h.db.Create(&interaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create interaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": interaction})
}

func (h *InteractionHandler) UpdateInteraction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interaction ID"})
		return
	}

	var req UpdateInteractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var interaction models.InteractionLog
	if err := h.db.First(&interaction, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interaction not found"})
		return
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if req.ContactID != nil {
		updates["contact_id"] = req.ContactID
	}
	if req.InteractionType != "" {
		updates["interaction_type"] = req.InteractionType
	}
	if req.InteractionDate != "" {
		interactionDate, _ := time.Parse("2006-01-02", req.InteractionDate)
		updates["interaction_date"] = interactionDate
	}
	if req.Topic != "" {
		updates["topic"] = req.Topic
	}
	if req.ContentSummary != "" {
		updates["content_summary"] = req.ContentSummary
	}
	if req.Outcome != "" {
		updates["outcome"] = req.Outcome
	}
	if req.NextAction != "" {
		updates["next_action"] = req.NextAction
	}
	if req.NextActionDate != nil {
		nextDate, _ := time.Parse("2006-01-02", *req.NextActionDate)
		updates["next_action_date"] = nextDate
	}
	if req.ImportanceLevel > 0 {
		updates["importance_level"] = req.ImportanceLevel
	}

	if err := h.db.Model(&interaction).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update interaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interaction})
}

func (h *InteractionHandler) DeleteInteraction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid interaction ID"})
		return
	}

	if err := h.db.Delete(&models.InteractionLog{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete interaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Interaction deleted successfully"})
}

func (h *InteractionHandler) GetInteractionTimeline(c *gin.Context) {
	customerID := c.Query("customer_id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	var interactions []models.InteractionLog
	if err := h.db.Where("customer_id = ?", customerID).
		Order("interaction_date DESC").
		Find(&interactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch timeline"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": interactions})
}

func (h *InteractionHandler) GetInteractionStats(c *gin.Context) {
	customerID := c.Query("customer_id")

	type TypeCount struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}

	type MonthlyCount struct {
		Month string `json:"month"`
		Count int64  `json:"count"`
	}

	type ContactCount struct {
		ContactName string `json:"contact_name"`
		Count       int64  `json:"count"`
	}

	var stats struct {
		TotalCount    int64          `json:"total_count"`
		TypeBreakdown []TypeCount    `json:"type_breakdown"`
		MonthlyTrend  []MonthlyCount `json:"monthly_trend"`
		TopContacts   []ContactCount `json:"top_contacts"`
	}

	query := h.db.Model(&models.InteractionLog{})
	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	query.Count(&stats.TotalCount)

	h.db.Model(&models.InteractionLog{}).
		Select("interaction_type as type, count(*) as count").
		Group("interaction_type").
		Scan(&stats.TypeBreakdown)

	h.db.Model(&models.InteractionLog{}).
		Select("TO_CHAR(interaction_date, 'YYYY-MM') as month, count(*) as count").
		Where("interaction_date > ?", time.Now().AddDate(-1, 0, 0)).
		Group("TO_CHAR(interaction_date, 'YYYY-MM')").
		Order("month DESC").
		Scan(&stats.MonthlyTrend)

	c.JSON(http.StatusOK, gin.H{"data": stats})
}
