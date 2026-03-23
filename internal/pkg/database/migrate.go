package database

import (
	"log"

	"enterprise-orbit/internal/models"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database migration...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.Customer{},
		&models.UserCustomer{},
		&models.Contact{},
		&models.OrgStructure{},
		&models.CooperationHistory{},
		&models.RelationshipHealth{},
		&models.ContractTemplate{},
		&models.Contract{},
		&models.ContractFile{},
		&models.ContractContent{},
		&models.ContractApprovalFlow{},
		&models.ContractPerformance{},
		&models.ContractRenewalPrediction{},
		&models.InteractionLog{},
		&models.InteractionAttachment{},
		&models.InteractionTag{},
		&models.KnowledgeExtract{},
	)
	if err != nil {
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}

func InitRBACData(db *gorm.DB) error {
	log.Println("Initializing RBAC data...")

	permissions := []models.Permission{
		{Resource: "user", Action: "create", Description: "Create users"},
		{Resource: "user", Action: "read", Description: "Read users"},
		{Resource: "user", Action: "update", Description: "Update users"},
		{Resource: "user", Action: "delete", Description: "Delete users"},
		{Resource: "client", Action: "create", Description: "Create clients"},
		{Resource: "client", Action: "read", Description: "Read clients"},
		{Resource: "client", Action: "update", Description: "Update clients"},
		{Resource: "client", Action: "delete", Description: "Delete clients"},
		{Resource: "role", Action: "create", Description: "Create roles"},
		{Resource: "role", Action: "read", Description: "Read roles"},
		{Resource: "role", Action: "update", Description: "Update roles"},
		{Resource: "role", Action: "delete", Description: "Delete roles"},
		{Resource: "permission", Action: "create", Description: "Create permissions"},
		{Resource: "permission", Action: "read", Description: "Read permissions"},
		{Resource: "permission", Action: "update", Description: "Update permissions"},
		{Resource: "permission", Action: "delete", Description: "Delete permissions"},
		{Resource: "customer", Action: "create", Description: "Create customers"},
		{Resource: "customer", Action: "read", Description: "Read customers"},
		{Resource: "customer", Action: "update", Description: "Update customers"},
		{Resource: "customer", Action: "delete", Description: "Delete customers"},
		{Resource: "contract", Action: "create", Description: "Create contracts"},
		{Resource: "contract", Action: "read", Description: "Read contracts"},
		{Resource: "contract", Action: "update", Description: "Update contracts"},
		{Resource: "contract", Action: "delete", Description: "Delete contracts"},
		{Resource: "contract", Action: "approve", Description: "Approve contracts"},
		{Resource: "interaction", Action: "create", Description: "Create interactions"},
		{Resource: "interaction", Action: "read", Description: "Read interactions"},
		{Resource: "interaction", Action: "update", Description: "Update interactions"},
		{Resource: "interaction", Action: "delete", Description: "Delete interactions"},
		{Resource: "dashboard", Action: "read", Description: "View dashboard"},
	}

	for _, perm := range permissions {
		var existing models.Permission
		result := db.Where("resource = ? AND action = ?", perm.Resource, perm.Action).First(&existing)
		if result.Error != nil {
			if err := db.Create(&perm).Error; err != nil {
				log.Printf("Failed to create permission %s:%s: %v", perm.Resource, perm.Action, err)
			}
		}
	}

	roles := []models.Role{
		{Name: "admin", Description: "Administrator with all permissions", Type: "static"},
		{Name: "user", Description: "Regular user", Type: "static"},
		{Name: "viewer", Description: "Read-only user", Type: "static"},
	}

	for _, role := range roles {
		var existing models.Role
		result := db.Where("name = ?", role.Name).First(&existing)
		if result.Error != nil {
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Failed to create role %s: %v", role.Name, err)
			}
		}
	}

	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err == nil {
		var adminPerms []models.Permission
		db.Find(&adminPerms)

		for _, perm := range adminPerms {
			var existing models.RolePermission
			result := db.Where("role_id = ? AND permission_id = ?", adminRole.ID, perm.ID).First(&existing)
			if result.Error != nil {
				rolePerm := models.RolePermission{
					RoleID:       adminRole.ID,
					PermissionID: perm.ID,
				}
				db.Create(&rolePerm)
			}
		}
	}

	var userRole models.Role
	if err := db.Where("name = ?", "user").First(&userRole).Error; err == nil {
		var readPerms []models.Permission
		db.Where("action = ?", "read").Find(&readPerms)

		for _, perm := range readPerms {
			var existing models.RolePermission
			result := db.Where("role_id = ? AND permission_id = ?", userRole.ID, perm.ID).First(&existing)
			if result.Error != nil {
				rolePerm := models.RolePermission{
					RoleID:       userRole.ID,
					PermissionID: perm.ID,
				}
				db.Create(&rolePerm)
			}
		}
	}

	log.Println("RBAC data initialization completed")
	return nil
}
