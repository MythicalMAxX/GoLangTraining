package migrations

import (
	"mypackage/internal/models"

	"gorm.io/gorm"
)

// RunMigrations performs all database migrations
func RunMigrations(db *gorm.DB) error {
	// Create enum types if they don't exist
	db.Exec(`DO $$ BEGIN
        CREATE TYPE membership_type AS ENUM ('standard', 'premium', 'vip');
        EXCEPTION WHEN duplicate_object THEN NULL;
    END $$;`)

	db.Exec(`DO $$ BEGIN
        CREATE TYPE member_status AS ENUM ('active', 'inactive', 'suspended');
        EXCEPTION WHEN duplicate_object THEN NULL;
    END $$;`)

	// Run auto migrations
	return db.AutoMigrate(&models.Member{}, &models.Borrow{})
}
