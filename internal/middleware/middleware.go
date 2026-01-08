package middleware

import (
    "net/http"
    "strings"
    "time"
    
    "enterprise-orbit/internal/models"
    
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// RBAC中间件
func RBAC(requiredPermission string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取当前用户
        user, exists := c.Get("currentUser")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "用户未认证"})
            c.Abort()
            return
        }
        
        currentUser := user.(*models.User)
        
        // 2. 检查是否为超级管理员
        if isSuperAdmin(currentUser) {
            c.Next()
            return
        }
        
        // 3. 动态角色评估
        err := evaluateDynamicRoles(c, currentUser)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "权限评估失败"})
            c.Abort()
            return
        }
        
        // 4. 权限检查
        hasPermission, err := checkPermission(c, currentUser, requiredPermission)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
            c.Abort()
            return
        }
        
        if !hasPermission {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "权限不足",
                "required_permission": requiredPermission,
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// 检查权限的核心逻辑
func checkPermission(c *gin.Context, user *models.User, requiredPermission string) (bool, error) {
    db := c.MustGet("db").(*gorm.DB)
    
    // 解析所需权限 (格式: "resource:action")
    parts := strings.Split(requiredPermission, ":")
    if len(parts) != 2 {
        return false, nil
    }
    resource, action := parts[0], parts[1]
    
    // 查询用户有效权限
    var permissionCount int64
    err := db.Table("permissions p").
        Joins("JOIN role_permissions rp ON rp.permission_id = p.id").
        Joins("JOIN roles r ON r.id = rp.role_id").
        Joins("JOIN user_roles ur ON ur.role_id = r.id").
        Where("ur.user_id = ? AND p.resource = ? AND p.action = ?", user.ID, resource, action).
        Where("ur.expires_at IS NULL OR ur.expires_at > ?", time.Now()).
        Count(&permissionCount).Error
        
    return permissionCount > 0, err
}

// 动态角色评估
func evaluateDynamicRoles(c *gin.Context, user *models.User) error {
    db := c.MustGet("db").(*gorm.DB)
    
    // 获取动态角色定义
    var dynamicRoles []models.Role
    if err := db.Where("type = ?", "dynamic").Find(&dynamicRoles).Error; err != nil {
        return err
    }
    
    for _, role := range dynamicRoles {
        // 评估动态角色条件
        if shouldAssignDynamicRole(user, role.Conditions) {
            expiresAt := time.Now().Add(24 * time.Hour)
            userRole := models.UserRole{
                UserID:    user.ID,
                RoleID:    role.ID,
                ExpiresAt: &expiresAt,
            }
            
            // 避免重复分配
            var existingCount int64
            db.Model(&models.UserRole{}).
                Where("user_id = ? AND role_id = ?", user.ID, role.ID).
                Count(&existingCount)
                
            if existingCount == 0 {
                db.Create(&userRole)
            }
        }
    }
    
    return nil
}

func isSuperAdmin(user *models.User) bool {
    return user.Username == "admin"
}

func shouldAssignDynamicRole(user *models.User, conditions string) bool {
    return conditions != "" && strings.Contains(conditions, user.Username)
}