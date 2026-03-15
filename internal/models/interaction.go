package models

import "time"

type InteractionLog struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	CustomerID       uint       `gorm:"not null;index" json:"customer_id"`
	ContactID        *uint      `json:"contact_id"`
	InteractionType  string     `gorm:"size:50;not null" json:"interaction_type"`
	InteractionDate  time.Time  `json:"interaction_date"`
	Topic            string     `gorm:"size:200" json:"topic"`
	ContentSummary   string     `gorm:"type:text" json:"content_summary"`
	Outcome          string     `gorm:"size:200" json:"outcome"`
	NextAction       string     `gorm:"size:200" json:"next_action"`
	NextActionDate   *time.Time `json:"next_action_date"`
	SentimentScore   float64    `gorm:"default:0" json:"sentiment_score"`
	ImportanceLevel  int        `gorm:"default:1" json:"importance_level"`
	CreatedBy        uint       `json:"created_by"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	Attachments []InteractionAttachment `gorm:"foreignKey:InteractionID" json:"attachments,omitempty"`
	Tags        []InteractionTag        `gorm:"many2many:interaction_log_tags;" json:"tags,omitempty"`
}

type InteractionAttachment struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	InteractionID  uint      `gorm:"not null;index" json:"interaction_id"`
	FileName       string    `gorm:"size:200;not null" json:"file_name"`
	FilePath       string    `gorm:"size:500;not null" json:"file_path"`
	FileSize       int64     `json:"file_size"`
	FileType       string    `gorm:"size:50" json:"file_type"`
	UploadedBy     uint      `json:"uploaded_by"`
	UploadedAt     time.Time `json:"uploaded_at"`
}

type InteractionTag struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TagName     string    `gorm:"size:50;not null" json:"tag_name"`
	TagType     string    `gorm:"size:50" json:"tag_type"`
	Description string    `gorm:"size:200" json:"description"`
	Color       string    `gorm:"size:20" json:"color"`
	CreatedAt   time.Time `json:"created_at"`
}

type KnowledgeExtract struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	InteractionID   uint      `gorm:"not null;index" json:"interaction_id"`
	KnowledgeType   string    `gorm:"size:50" json:"knowledge_type"`
	KnowledgeTitle  string    `gorm:"size:200" json:"knowledge_title"`
	KnowledgeContent string   `gorm:"type:text" json:"knowledge_content"`
	ConfidenceScore float64   `gorm:"default:0" json:"confidence_score"`
	ExtractedBy     uint      `json:"extracted_by"`
	ExtractedAt     time.Time `json:"extracted_at"`
}
