package models

import (
	"gorm.io/gorm"
	"time"
)

type Client struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CompanyName   string         `gorm:"size:200;not null" json:"company_name"`
	Industry      string         `gorm:"size:100" json:"industry"`
	ContactPerson string         `gorm:"size:100" json:"contact_person"`
	Email         string         `gorm:"size:100" json:"email"`
	Phone         string         `gorm:"size:20" json:"phone"`
	Address       string         `gorm:"type:text" json:"address"`
	Status        string         `gorm:"size:20;default:potential" json:"status"` // potential, active, lost
	Level         string         `gorm:"size:20;default:normal" json:"level"`     // important, normal, trial
	Tags          string         `gorm:"type:text" json:"tags"`                   // JSON string or comma separated
	CreatedBy     uint           `gorm:"not null" json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	
	// Associations
	Contracts []Contract `gorm:"foreignKey:ClientID" json:"contracts,omitempty"`
	Interactions []Interaction `gorm:"foreignKey:ClientID" json:"interactions,omitempty"`
}