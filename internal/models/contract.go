package models

import "time"

type ContractTemplate struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TemplateName string    `gorm:"size:200;not null" json:"template_name"`
	TemplateType string    `gorm:"size:50" json:"template_type"`
	Content      string    `gorm:"type:text" json:"content"`
	Variables    string    `gorm:"type:text" json:"variables"`
	Status       string    `gorm:"size:20;default:active" json:"status"`
	CreatedBy    uint      `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Contract struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	ContractNo      string     `gorm:"uniqueIndex;size:50;not null" json:"contract_no"`
	ContractName    string     `gorm:"size:200;not null" json:"contract_name"`
	CustomerID      uint       `gorm:"not null;index" json:"customer_id"`
	TemplateID      *uint      `json:"template_id"`
	ContractType    string     `gorm:"size:50" json:"contract_type"`
	Amount          float64    `gorm:"default:0" json:"amount"`
	Currency        string     `gorm:"size:10;default:CNY" json:"currency"`
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	Status          string     `gorm:"size:20;default:draft" json:"status"`
	CurrentApprover *uint      `json:"current_approver"`
	ApprovalStatus  string     `gorm:"size:20;default:pending" json:"approval_status"`
	CreatedBy       uint       `json:"created_by"`
	SignDate        *time.Time `json:"sign_date"`
	EffectiveDate   *time.Time `json:"effective_date"`
	ExpiryAlertSent bool       `gorm:"default:false" json:"expiry_alert_sent"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	Files []ContractFile `gorm:"foreignKey:ContractID" json:"files,omitempty"`
}

type ContractFile struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ContractID   uint      `gorm:"not null;index" json:"contract_id"`
	FileName     string    `gorm:"size:255;not null" json:"file_name"`
	FilePath     string    `gorm:"size:500;not null" json:"file_path"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `gorm:"size:50" json:"file_type"`
	Category     string    `gorm:"size:50;default:contract" json:"category"`
	Description  string    `gorm:"size:500" json:"description"`
	UploadedBy   uint      `json:"uploaded_by"`
	IsShared     bool      `gorm:"default:false" json:"is_shared"`
	CreatedAt    time.Time `json:"created_at"`

	Contract Contract `gorm:"foreignKey:ContractID" json:"contract,omitempty"`
}

type ContractContent struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ContractID  uint      `gorm:"not null;index" json:"contract_id"`
	ContentJSON string    `gorm:"type:text" json:"content_json"`
	Version     int       `gorm:"default:1" json:"version"`
	ChangeLog   string    `gorm:"type:text" json:"change_log"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

type ContractApprovalFlow struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ContractID     uint       `gorm:"not null;index" json:"contract_id"`
	ApproverID     uint       `gorm:"not null" json:"approver_id"`
	ApprovalOrder  int        `gorm:"default:1" json:"approval_order"`
	ApprovalStatus string     `gorm:"size:20;default:pending" json:"approval_status"`
	Comment        string     `gorm:"type:text" json:"comment"`
	ApprovalTime   *time.Time `json:"approval_time"`
	CreatedAt      time.Time  `json:"created_at"`
}

type ContractPerformance struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	ContractID       uint       `gorm:"not null;index" json:"contract_id"`
	MilestoneName    string     `gorm:"size:200;not null" json:"milestone_name"`
	MilestoneDesc    string     `gorm:"type:text" json:"milestone_desc"`
	PlannedDate      time.Time  `json:"planned_date"`
	ActualDate       *time.Time `json:"actual_date"`
	CompletionStatus string     `gorm:"size:20;default:pending" json:"completion_status"`
	AmountPlanned    float64    `gorm:"default:0" json:"amount_planned"`
	AmountPaid       float64    `gorm:"default:0" json:"amount_paid"`
	Notes            string     `gorm:"type:text" json:"notes"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type ContractRenewalPrediction struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	ContractID             uint      `gorm:"uniqueIndex;not null" json:"contract_id"`
	RenewalProbability     float64   `gorm:"default:0" json:"renewal_probability"`
	PredictedRenewalAmount float64   `gorm:"default:0" json:"predicted_renewal_amount"`
	ConfidenceLevel        float64   `gorm:"default:0" json:"confidence_level"`
	PredictionFactors      string    `gorm:"type:text" json:"prediction_factors"`
	RiskFactors            string    `gorm:"type:text" json:"risk_factors"`
	Recommendation         string    `gorm:"type:text" json:"recommendation"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}
