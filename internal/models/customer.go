package models

import "time"

type Customer struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CustomerCode   string    `gorm:"uniqueIndex;size:50;not null" json:"customer_code"`
	CustomerName   string    `gorm:"size:200;not null" json:"customer_name"`
	Industry       string    `gorm:"size:100" json:"industry"`
	Scale          string    `gorm:"size:50" json:"scale"`
	Region         string    `gorm:"size:100" json:"region"`
	Address        string    `gorm:"size:500" json:"address"`
	Website        string    `gorm:"size:200" json:"website"`
	Status         string    `gorm:"size:20;default:active" json:"status"`
	TrustScore     float64   `gorm:"default:0" json:"trust_score"`
	CommitmentScore float64  `gorm:"default:0" json:"commitment_score"`
	ReciprocityScore float64 `gorm:"default:0" json:"reciprocity_score"`
	OverallScore   float64   `gorm:"default:0" json:"overall_score"`
	CreatedBy      uint      `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Contacts           []Contact           `gorm:"foreignKey:CustomerID" json:"contacts,omitempty"`
	OrgStructures      []OrgStructure      `gorm:"foreignKey:CustomerID" json:"org_structures,omitempty"`
	CooperationHistory []CooperationHistory `gorm:"foreignKey:CustomerID" json:"cooperation_history,omitempty"`
	Contracts          []Contract          `gorm:"foreignKey:CustomerID" json:"contracts,omitempty"`
	InteractionLogs    []InteractionLog    `gorm:"foreignKey:CustomerID" json:"interaction_logs,omitempty"`
}

type Contact struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	CustomerID        uint      `gorm:"not null;index" json:"customer_id"`
	Name              string    `gorm:"size:100;not null" json:"name"`
	Position          string    `gorm:"size:100" json:"position"`
	Department        string    `gorm:"size:100" json:"department"`
	Phone             string    `gorm:"size:50" json:"phone"`
	Email             string    `gorm:"size:100" json:"email"`
	Wechat            string    `gorm:"size:50" json:"wechat"`
	IsKeyDecisionMaker bool     `gorm:"default:false" json:"is_key_decision_maker"`
	InfluenceLevel    int       `gorm:"default:1" json:"influence_level"`
	Notes             string    `gorm:"type:text" json:"notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	Customer Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

type OrgStructure struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CustomerID   uint      `gorm:"not null;index" json:"customer_id"`
	ParentID     *uint     `json:"parent_id"`
	DeptName     string    `gorm:"size:100;not null" json:"dept_name"`
	DeptFunction string    `gorm:"size:200" json:"dept_function"`
	OrderNum     int       `gorm:"default:0" json:"order_num"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Customer Customer      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Parent   *OrgStructure `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []OrgStructure `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

type CooperationHistory struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	CustomerID       uint      `gorm:"not null;index" json:"customer_id"`
	ProjectName      string    `gorm:"size:200;not null" json:"project_name"`
	ContractAmount   float64   `gorm:"default:0" json:"contract_amount"`
	StartDate        time.Time `json:"start_date"`
	EndDate          *time.Time `json:"end_date"`
	Status           string    `gorm:"size:20" json:"status"`
	SatisfactionScore float64  `gorm:"default:0" json:"satisfaction_score"`
	Description      string    `gorm:"type:text" json:"description"`
	CreatedAt        time.Time `json:"created_at"`

	Customer Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

type RelationshipHealth struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	CustomerID          uint      `gorm:"uniqueIndex;not null" json:"customer_id"`
	TrustScore          float64   `gorm:"default:0" json:"trust_score"`
	CommitmentScore     float64   `gorm:"default:0" json:"commitment_score"`
	ReciprocityScore    float64   `gorm:"default:0" json:"reciprocity_score"`
	OverallScore        float64   `gorm:"default:0" json:"overall_score"`
	InteractionBalance  float64   `gorm:"default:0" json:"interaction_balance"`
	ResponseTime        float64   `gorm:"default:0" json:"response_time"`
	AssessmentDate      time.Time `json:"assessment_date"`
	NextAssessmentDate  time.Time `json:"next_assessment_date"`
	RiskLevel           string    `gorm:"size:20;default:low" json:"risk_level"`
	Recommendations     string    `gorm:"type:text" json:"recommendations"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	Customer Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}
