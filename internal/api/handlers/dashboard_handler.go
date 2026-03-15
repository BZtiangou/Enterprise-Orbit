package handlers

import (
	"net/http"
	"time"

	"enterprise-orbit/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) GetOverview(c *gin.Context) {
	var overview struct {
		TotalCustomers      int64   `json:"total_customers"`
		ActiveCustomers     int64   `json:"active_customers"`
		TotalContracts      int64   `json:"total_contracts"`
		ActiveContracts     int64   `json:"active_contracts"`
		TotalContractValue  float64 `json:"total_contract_value"`
		TotalInteractions   int64   `json:"total_interactions"`
		MonthlyInteractions int64   `json:"monthly_interactions"`
		ExpiringContracts   int64   `json:"expiring_contracts"`
		HighRiskCustomers   int64   `json:"high_risk_customers"`
	}

	h.db.Model(&models.Customer{}).Count(&overview.TotalCustomers)
	h.db.Model(&models.Customer{}).Where("status = ?", "active").Count(&overview.ActiveCustomers)
	h.db.Model(&models.Contract{}).Count(&overview.TotalContracts)
	h.db.Model(&models.Contract{}).Where("status = ?", "active").Count(&overview.ActiveContracts)
	h.db.Model(&models.Contract{}).Where("status = ?", "active").Select("COALESCE(SUM(amount), 0)").Scan(&overview.TotalContractValue)
	h.db.Model(&models.InteractionLog{}).Count(&overview.TotalInteractions)
	h.db.Model(&models.InteractionLog{}).Where("interaction_date >= ?", time.Now().AddDate(0, -1, 0)).Count(&overview.MonthlyInteractions)
	h.db.Model(&models.Contract{}).Where("status = ? AND end_date <= ? AND end_date > ?", "active", time.Now().AddDate(0, 0, 30), time.Now()).Count(&overview.ExpiringContracts)
	h.db.Model(&models.RelationshipHealth{}).Where("risk_level = ?", "high").Count(&overview.HighRiskCustomers)

	c.JSON(http.StatusOK, gin.H{"data": overview})
}

func (h *DashboardHandler) GetCustomerGrowthTrend(c *gin.Context) {
	months := c.DefaultQuery("months", "12")
	monthsInt := 12
	if months != "" {
		if m, err := time.ParseDuration(months + "h"); err == nil {
			monthsInt = int(m.Hours() / 720)
		}
	}

	var trend []struct {
		Month  string  `json:"month"`
		New    int64   `json:"new"`
		Total  int64   `json:"total"`
		Growth float64 `json:"growth"`
	}

	startDate := time.Now().AddDate(0, -monthsInt, 0)

	h.db.Model(&models.Customer{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, count(*) as new").
		Where("created_at >= ?", startDate).
		Group("TO_CHAR(created_at, 'YYYY-MM')").
		Order("month ASC").
		Scan(&trend)

	var runningTotal int64
	for i := range trend {
		runningTotal += trend[i].New
		trend[i].Total = runningTotal
		if i > 0 {
			if trend[i-1].Total > 0 {
				trend[i].Growth = float64(trend[i].New) / float64(trend[i-1].Total) * 100
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": trend})
}

func (h *DashboardHandler) GetContractValueDistribution(c *gin.Context) {
	var distribution []struct {
		Status string  `json:"status"`
		Count  int64   `json:"count"`
		Amount float64 `json:"amount"`
	}

	h.db.Model(&models.Contract{}).
		Select("status, count(*) as count, sum(amount) as amount").
		Group("status").
		Scan(&distribution)

	c.JSON(http.StatusOK, gin.H{"data": distribution})
}

func (h *DashboardHandler) GetInteractionHeatmap(c *gin.Context) {
	year := c.DefaultQuery("year", time.Now().Format("2006"))

	var heatmap []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}

	h.db.Model(&models.InteractionLog{}).
		Select("TO_CHAR(interaction_date, 'YYYY-MM-DD') as date, count(*) as count").
		Where("TO_CHAR(interaction_date, 'YYYY') = ?", year).
		Group("TO_CHAR(interaction_date, 'YYYY-MM-DD')").
		Order("date ASC").
		Scan(&heatmap)

	c.JSON(http.StatusOK, gin.H{"data": heatmap})
}

func (h *DashboardHandler) GetRenewalRiskAlerts(c *gin.Context) {
	var alerts []struct {
		ContractID         uint    `json:"contract_id"`
		ContractNo         string  `json:"contract_no"`
		ContractName       string  `json:"contract_name"`
		CustomerID         uint    `json:"customer_id"`
		CustomerName       string  `json:"customer_name"`
		EndDate            string  `json:"end_date"`
		Amount             float64 `json:"amount"`
		DaysUntilExpiry    int     `json:"days_until_expiry"`
		RenewalProbability float64 `json:"renewal_probability"`
		RiskLevel          string  `json:"risk_level"`
	}

	h.db.Table("contracts c").
		Select(`c.id as contract_id, c.contract_no, c.contract_name, 
				c.customer_id, cu.customer_name, c.end_date, c.amount,
				DATE_PART('day', c.end_date - CURRENT_DATE) as days_until_expiry,
				COALESCE(cr.renewal_probability, 50) as renewal_probability,
				CASE 
					WHEN COALESCE(cr.renewal_probability, 50) >= 70 THEN 'low'
					WHEN COALESCE(cr.renewal_probability, 50) >= 40 THEN 'medium'
					ELSE 'high'
				END as risk_level`).
		Joins("LEFT JOIN customers cu ON c.customer_id = cu.id").
		Joins("LEFT JOIN contract_renewal_predictions cr ON c.id = cr.contract_id").
		Where("c.status = ? AND c.end_date <= ? AND c.end_date > ?", "active", time.Now().AddDate(0, 0, 90), time.Now()).
		Order("days_until_expiry ASC").
		Limit(10).
		Scan(&alerts)

	c.JSON(http.StatusOK, gin.H{"data": alerts})
}

func (h *DashboardHandler) GetTopCustomers(c *gin.Context) {
	var customers []struct {
		CustomerID     uint    `json:"customer_id"`
		CustomerName   string  `json:"customer_name"`
		TotalContracts int64   `json:"total_contracts"`
		TotalValue     float64 `json:"total_value"`
		HealthScore    float64 `json:"health_score"`
	}

	h.db.Table("customers cu").
		Select(`cu.id as customer_id, cu.customer_name, 
				count(c.id) as total_contracts, 
				COALESCE(sum(c.amount), 0) as total_value,
				COALESCE(rh.overall_score, 0) as health_score`).
		Joins("LEFT JOIN contracts c ON cu.id = c.customer_id AND c.status = 'active'").
		Joins("LEFT JOIN relationship_health rh ON cu.id = rh.customer_id").
		Where("cu.status = ?", "active").
		Group("cu.id, cu.customer_name, rh.overall_score").
		Order("total_value DESC").
		Limit(10).
		Scan(&customers)

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

func (h *DashboardHandler) GetRecentActivities(c *gin.Context) {
	var activities []struct {
		Type        string    `json:"type"`
		ID          uint      `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
	}

	var interactions []models.InteractionLog
	h.db.Order("created_at DESC").Limit(10).Find(&interactions)
	for _, i := range interactions {
		activities = append(activities, struct {
			Type        string    `json:"type"`
			ID          uint      `json:"id"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			CreatedAt   time.Time `json:"created_at"`
		}{
			Type:        "interaction",
			ID:          i.ID,
			Title:       i.Topic,
			Description: i.ContentSummary,
			CreatedAt:   i.CreatedAt,
		})
	}

	var contracts []models.Contract
	h.db.Order("created_at DESC").Limit(10).Find(&contracts)
	for _, c := range contracts {
		activities = append(activities, struct {
			Type        string    `json:"type"`
			ID          uint      `json:"id"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			CreatedAt   time.Time `json:"created_at"`
		}{
			Type:        "contract",
			ID:          c.ID,
			Title:       c.ContractName,
			Description: "Contract " + c.Status,
			CreatedAt:   c.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": activities})
}
