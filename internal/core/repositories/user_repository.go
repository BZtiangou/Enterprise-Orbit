package repositories

import (
	"enterprise-orbit/internal/models"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	HasPermission(userID uint, resource, action string) (bool, error)
	GetUserRoles(userID uint) ([]models.Role, error)
	AssignRole(userID, roleID uint, expiresAt *time.Time) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) HasPermission(userID uint, resource, action string) (bool, error) {
	var count int64
	r.db.Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Joins("JOIN role_permissions ON roles.id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ?", userID, resource, action).
		Where("user_roles.expires_at IS NULL OR user_roles.expires_at > ?", time.Now()).
		Count(&count)
	return count > 0, nil
}

func (r *userRepository) GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Table("user_roles").
		Joins("JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Where("user_roles.expires_at IS NULL OR user_roles.expires_at > ?", time.Now()).
		Find(&roles).Error
	return roles, err
}

func (r *userRepository) AssignRole(userID, roleID uint, expiresAt *time.Time) error {
	userRole := models.UserRole{
		UserID:    userID,
		RoleID:    roleID,
		ExpiresAt: expiresAt,
	}
	return r.db.Table("user_roles").Create(&userRole).Error
}
