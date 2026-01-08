package models

import (
	"time"
)

type Client struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Code      string    `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Type      string    `gorm:"size:20;default:enterprise" json:"type"` // enterprise, individual
	Status    string    `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
