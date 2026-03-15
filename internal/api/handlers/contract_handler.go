package handlers

import (
	"net/http"
	"strconv"
	"time"

	"enterprise-orbit/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContractHandler struct {
	db *gorm.DB
}

func NewContractHandler(db *gorm.DB) *ContractHandler {
	return &ContractHandler{db: db}
}

type CreateContractRequest struct {
	ContractNo   string  `json:"contract_no" binding:"required"`
	ContractName string  `json:"contract_name" binding:"required"`
	CustomerID   uint    `json:"customer_id" binding:"required"`
	TemplateID   *uint   `json:"template_id"`
	ContractType string  `json:"contract_type"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	StartDate    string  `json:"start_date" binding:"required"`
	EndDate      string  `json:"end_date" binding:"required"`
}

type UpdateContractRequest struct {
	ContractName string  `json:"contract_name"`
	ContractType string  `json:"contract_type"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	Status       string  `json:"status"`
}

type ApprovalRequest struct {
	Status  string `json:"status" binding:"required"`
	Comment string `json:"comment"`
}

func (h *ContractHandler) GetContracts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	status := c.Query("status")
	customerID := c.Query("customer_id")

	var contracts []models.Contract
	var total int64

	query := h.db.Model(&models.Contract{}).Preload("Customer")

	if search != "" {
		query = query.Where("contract_name LIKE ? OR contract_no LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&contracts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contracts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": contracts,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

func (h *ContractHandler) GetContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var contract models.Contract
	if err := h.db.Preload("Customer").
		Preload("Contents").
		Preload("ApprovalFlows").
		Preload("PerformanceRecords").
		First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": contract})
}

func (h *ContractHandler) CreateContract(c *gin.Context) {
	var req CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.GetUint("userID")

	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	contract := models.Contract{
		ContractNo:     req.ContractNo,
		ContractName:   req.ContractName,
		CustomerID:     req.CustomerID,
		TemplateID:     req.TemplateID,
		ContractType:   req.ContractType,
		Amount:         req.Amount,
		Currency:       req.Currency,
		StartDate:      startDate,
		EndDate:        endDate,
		Status:         "draft",
		ApprovalStatus: "pending",
		CreatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.db.Create(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contract"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": contract})
}

func (h *ContractHandler) UpdateContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var req UpdateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var contract models.Contract
	if err := h.db.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if req.ContractName != "" {
		updates["contract_name"] = req.ContractName
	}
	if req.ContractType != "" {
		updates["contract_type"] = req.ContractType
	}
	if req.Amount > 0 {
		updates["amount"] = req.Amount
	}
	if req.Currency != "" {
		updates["currency"] = req.Currency
	}
	if req.StartDate != "" {
		startDate, _ := time.Parse("2006-01-02", req.StartDate)
		updates["start_date"] = startDate
	}
	if req.EndDate != "" {
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		updates["end_date"] = endDate
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := h.db.Model(&contract).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contract"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": contract})
}

func (h *ContractHandler) DeleteContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	if err := h.db.Delete(&models.Contract{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contract"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contract deleted successfully"})
}

func (h *ContractHandler) SubmitForApproval(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var contract models.Contract
	if err := h.db.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	if contract.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only draft contracts can be submitted"})
		return
	}

	contract.Status = "pending_approval"
	contract.ApprovalStatus = "pending"

	if err := h.db.Save(&contract).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit contract"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contract submitted for approval", "data": contract})
}

func (h *ContractHandler) ApproveContract(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var req ApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.GetUint("userID")

	var contract models.Contract
	if err := h.db.First(&contract, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	approvalFlow := models.ContractApprovalFlow{
		ContractID:     contract.ID,
		ApproverID:     userID,
		ApprovalStatus: req.Status,
		Comment:        req.Comment,
		ApprovalTime:   &[]time.Time{time.Now()}[0],
		CreatedAt:      time.Now(),
	}

	if err := h.db.Create(&approvalFlow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record approval"})
		return
	}

	if req.Status == "approved" {
		contract.Status = "active"
		contract.ApprovalStatus = "approved"
		now := time.Now()
		contract.SignDate = &now
		contract.EffectiveDate = &now
	} else {
		contract.Status = "rejected"
		contract.ApprovalStatus = "rejected"
	}

	h.db.Save(&contract)

	c.JSON(http.StatusOK, gin.H{"message": "Approval recorded", "data": contract})
}

func (h *ContractHandler) GetExpiringContracts(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	var contracts []models.Contract
	expiryDate := time.Now().AddDate(0, 0, days)

	if err := h.db.Preload("Customer").
		Where("status = ? AND end_date <= ? AND end_date > ?", "active", expiryDate, time.Now()).
		Order("end_date ASC").
		Find(&contracts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contracts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": contracts})
}

func (h *ContractHandler) GetContractPerformance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var records []models.ContractPerformance
	if err := h.db.Where("contract_id = ?", id).Order("planned_date ASC").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch performance records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": records})
}

func (h *ContractHandler) CreatePerformanceRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	var req struct {
		MilestoneName string  `json:"milestone_name" binding:"required"`
		MilestoneDesc string  `json:"milestone_desc"`
		PlannedDate   string  `json:"planned_date" binding:"required"`
		AmountPlanned float64 `json:"amount_planned"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	plannedDate, _ := time.Parse("2006-01-02", req.PlannedDate)

	record := models.ContractPerformance{
		ContractID:       uint(id),
		MilestoneName:    req.MilestoneName,
		MilestoneDesc:    req.MilestoneDesc,
		PlannedDate:      plannedDate,
		CompletionStatus: "pending",
		AmountPlanned:    req.AmountPlanned,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := h.db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create performance record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": record})
}
