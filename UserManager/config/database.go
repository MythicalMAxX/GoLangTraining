package config

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
    dbURI := os.Getenv("USER_DB_URI")
    log.Printf("Connecting to database...")
    
    db, err := sql.Open("postgres", dbURI)
    if err != nil {
        log.Printf("Database connection error: %v", err)
        return nil, fmt.Errorf("error connecting to database: %v", err)
    }

    if err = db.Ping(); err != nil {
        log.Printf("Database ping error: %v", err)
        return nil, fmt.Errorf("error pinging database: %v", err)
    }

    log.Printf("Successfully connected to database")
    return db, nil
}

func runMigrations(db *sql.DB) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("could not create the postgres driver: %v", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations", // migrations files path
        "postgres",         // database name
        driver,            // database driver
    )
    if err != nil {
        return fmt.Errorf("error creating migration instance: %v", err)
    }

    // Run migrations
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("error running migrations: %v", err)
    }

    log.Println("Database migrations completed successfully")
    return nil
}