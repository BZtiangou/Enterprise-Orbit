package services

import (
    "enterprise-orbit/internal/models"
    "enterprise-orbit/internal/core/repositories"
    "time"
    
    "gorm.io/gorm"
)

type RBACService struct {
    userRepo repositories.UserRepository
    roleRepo repositories.RoleRepository
    permRepo repositories.PermissionRepository
    db       *gorm.DB
}

func NewRBACService(db *gorm.DB) *RBACService {
    return &RBACService{
        userRepo: repositories.NewUserRepository(db),
        roleRepo: repositories.NewRoleRepository(db),
        permRepo: repositories.NewPermissionRepository(db),
        db:       db,
    }
}

// 检查用户权限
func (s *RBACService) CheckPermission(userID uint, resource, action string) (bool, error) {
    // 实现权限检查逻辑
    return s.userRepo.HasPermission(userID, resource, action)
}

// 获取用户所有角色
func (s *RBACService) GetUserRoles(userID uint) ([]models.Role, error) {
    return s.userRepo.GetUserRoles(userID)
}

// 分配角色给用户
func (s *RBACService) AssignRoleToUser(userID, roleID uint, expiresAt *time.Time) error {
    return s.userRepo.AssignRole(userID, roleID, expiresAt)
}