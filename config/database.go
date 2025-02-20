package config

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "mypackage/migrations"
)

func InitDB() (*gorm.DB, error) {
    dsn := "postgresql://neondb_owner:npg_qXyMNZF4IK3n@ep-damp-heart-a8zcur1e-pooler.eastus2.azure.neon.tech/neondb?sslmode=require"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Run migrations
    if err := migrations.RunMigrations(db); err != nil {
        return nil, err
    }

    return db, nil
}