package models

import (
    "time"
    // "gorm.io/gorm"
)

// 角色模型
type Role struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
    Description string    `gorm:"size:200" json:"description"`
    Type        string    `gorm:"size:20;default:static" json:"type"` // static, dynamic
    Conditions  string    `gorm:"type:text" json:"conditions"` // 动态角色条件(JSON)
    CreatedAt   time.Time `json:"created_at"`
    
    // 关联
    UserRoles   []UserRole    `gorm:"foreignKey:RoleID" json:"user_roles,omitempty"`
    RolePermissions []RolePermission `gorm:"foreignKey:RoleID" json:"role_permissions,omitempty"`
}

// 权限模型
type Permission struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Resource    string    `gorm:"size:100;not null" json:"resource"`    // 资源: client, contract, user
    Action      string    `gorm:"size:50;not null" json:"action"`       // 操作: create, read, update, delete
    Description string    `gorm:"size:200" json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    
    // 关联
    RolePermissions []RolePermission `gorm:"foreignKey:PermissionID" json:"role_permissions,omitempty"`
}

// 用户-角色关联
type UserRole struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `gorm:"not null;index" json:"user_id"`
    RoleID    uint      `gorm:"not null;index" json:"role_id"`
    ExpiresAt *time.Time `json:"expires_at"` // 角色过期时间(动态特性)
    CreatedAt time.Time `json:"created_at"`
    
    // 关联
    User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Role Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// 角色-权限关联
type RolePermission struct {
    ID           uint       `gorm:"primaryKey" json:"id"`
    RoleID       uint       `gorm:"not null;index" json:"role_id"`
    PermissionID uint       `gorm:"not null;index" json:"permission_id"`
    Conditions   string     `gorm:"type:text" json:"conditions"` // 权限条件(JSON)
    CreatedAt    time.Time  `json:"created_at"`
    
    // 关联
    Role       Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
    Permission Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}