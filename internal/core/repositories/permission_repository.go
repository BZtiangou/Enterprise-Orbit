package repositories

import (
	"enterprise-orbit/internal/models"

	"gorm.io/gorm"
)

type PermissionRepository interface {
	GetAllPermissions() ([]models.Permission, error)
	GetPermissionsByRoleID(roleID uint) ([]models.Permission, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) GetAllPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) GetPermissionsByRoleID(roleID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Table("role_permissions").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error
	return permissions, err
}
