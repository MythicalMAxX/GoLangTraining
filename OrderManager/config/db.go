package config

import (
    "database/sql"
    "log"
    "os"
    "time"
    _ "github.com/lib/pq"
)

func RunMigrations(db *sql.DB) error {
    migrationSQL, err := os.ReadFile("migrations/001_create_orders_table.sql")
    if err != nil {
        return err
    }

    _, err = db.Exec(string(migrationSQL))
    if err != nil {
        return err
    }

    log.Println("Migrations completed successfully")
    return nil
}

func InitDB() *sql.DB {
    connectionString := os.Getenv("USER_DB_URI")
    if connectionString == "" {
        log.Fatal("USER_DB_URI environment variable is not set")
    }

    var db *sql.DB
    var err error
    
    // Retry connection up to 5 times
    for i := 0; i < 5; i++ {
        db, err = sql.Open("postgres", connectionString)
        if err != nil {
            log.Printf("Failed to open database connection: %v", err)
            time.Sleep(2 * time.Second)
            continue
        }

        err = db.Ping()
        if err == nil {
            log.Printf("Successfully connected to database")
            return db
        }
        
        log.Printf("Failed to ping database (attempt %d/5): %v", i+1, err)
        time.Sleep(2 * time.Second)
    }

    log.Fatalf("Could not establish database connection after 5 attempts: %v", err)
    return nil
}
