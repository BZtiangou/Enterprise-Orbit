package models

import (
    // "gorm.io/gorm"
    "time"
)

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
    Email     string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
    Password  string    `gorm:"not null" json:"-"`
    Status    string    `gorm:"size:20;default:active" json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    // 关联
    UserRoles []UserRole `gorm:"foreignKey:UserID" json:"user_roles,omitempty"`
}

type UserResponse struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
}