package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string    `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Nickname  string    `gorm:"size:50" json:"nickname"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	Department string   `gorm:"size:50" json:"department"`
	Position  string    `gorm:"size:50" json:"position"`
	Role      string    `gorm:"size:20;default:user" json:"role"`
	Status    string    `gorm:"size:20;default:active" json:"status"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserRoles []UserRole `gorm:"foreignKey:UserID" json:"user_roles,omitempty"`
}

type UserResponse struct {
	ID         uint       `json:"id"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Nickname   string     `json:"nickname"`
	Phone      string     `json:"phone"`
	Avatar     string     `json:"avatar"`
	Department string     `json:"department"`
	Position   string     `json:"position"`
	Role       string     `json:"role"`
	Status     string     `json:"status"`
	LastLogin  *time.Time `json:"last_login"`
	CreatedAt  time.Time  `json:"created_at"`
}

type UpdateProfileRequest struct {
	Nickname   string `json:"nickname"`
	Phone      string `json:"phone"`
	Avatar     string `json:"avatar"`
	Department string `json:"department"`
	Position   string `json:"position"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}
