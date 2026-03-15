package handlers

import (
	"net/http"
	"strconv"
	"time"

	"enterprise-orbit/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CustomerHandler struct {
	db *gorm.DB
}

func NewCustomerHandler(db *gorm.DB) *CustomerHandler {
	return &CustomerHandler{db: db}
}

type CreateCustomerRequest struct {
	CustomerCode string `json:"customer_code" binding:"required"`
	CustomerName string `json:"customer_name" binding:"required"`
	Industry     string `json:"industry"`
	Scale        string `json:"scale"`
	Region       string `json:"region"`
	Address      string `json:"address"`
	Website      string `json:"website"`
}

type UpdateCustomerRequest struct {
	CustomerName string `json:"customer_name"`
	Industry     string `json:"industry"`
	Scale        string `json:"scale"`
	Region       string `json:"region"`
	Address      string `json:"address"`
	Website      string `json:"website"`
	Status       string `json:"status"`
}

func (h *CustomerHandler) GetCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	status := c.Query("status")

	var customers []models.Customer
	var total int64

	query := h.db.Model(&models.Customer{})

	if search != "" {
		query = query.Where("customer_name LIKE ? OR customer_code LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": customers,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var customer models.Customer
	if err := h.db.Preload("Contacts").
		Preload("OrgStructures").
		Preload("CooperationHistory").
		First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.GetUint("userID")

	customer := models.Customer{
		CustomerCode: req.CustomerCode,
		CustomerName: req.CustomerName,
		Industry:     req.Industry,
		Scale:        req.Scale,
		Region:       req.Region,
		Address:      req.Address,
		Website:      req.Website,
		Status:       "active",
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.db.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": customer})
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var customer models.Customer
	if err := h.db.First(&customer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if req.CustomerName != "" {
		updates["customer_name"] = req.CustomerName
	}
	if req.Industry != "" {
		updates["industry"] = req.Industry
	}
	if req.Scale != "" {
		updates["scale"] = req.Scale
	}
	if req.Region != "" {
		updates["region"] = req.Region
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Website != "" {
		updates["website"] = req.Website
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := h.db.Model(&customer).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	if err := h.db.Delete(&models.Customer{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

func (h *CustomerHandler) GetCustomerHealth(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var health models.RelationshipHealth
	if err := h.db.Where("customer_id = ?", id).First(&health).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Health data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": health})
}

func (h *CustomerHandler) CalculateHealthScore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var interactionCount int64
	h.db.Model(&models.InteractionLog{}).Where("customer_id = ?", id).Count(&interactionCount)

	var contractCount int64
	h.db.Model(&models.Contract{}).Where("customer_id = ? AND status = ?", id, "active").Count(&contractCount)

	trustScore := float64(interactionCount) * 0.1
	if trustScore > 100 {
		trustScore = 100
	}

	commitmentScore := float64(contractCount) * 20
	if commitmentScore > 100 {
		commitmentScore = 100
	}

	reciprocityScore := (trustScore + commitmentScore) / 2
	overallScore := (trustScore + commitmentScore + reciprocityScore) / 3

	var health models.RelationshipHealth
	h.db.Where("customer_id = ?", id).FirstOrCreate(&health)

	health.CustomerID = uint(id)
	health.TrustScore = trustScore
	health.CommitmentScore = commitmentScore
	health.ReciprocityScore = reciprocityScore
	health.OverallScore = overallScore
	health.AssessmentDate = time.Now()
	health.NextAssessmentDate = time.Now().AddDate(0, 1, 0)

	if overallScore >= 70 {
		health.RiskLevel = "low"
	} else if overallScore >= 40 {
		health.RiskLevel = "medium"
	} else {
		health.RiskLevel = "high"
	}

	if err := h.db.Save(&health).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save health data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": health})
}
