package migrations

import (
	"mypackage/internal/models"

	"gorm.io/gorm"
)

// List of all models that need to be migrated
var Models = []interface{}{
	&models.Member{},
}

// RunMigrations performs all database migrations
func RunMigrations(db *gorm.DB) error {
	// This will automatically create/update tables based on struct definitions
	return db.AutoMigrate(Models...)
}
