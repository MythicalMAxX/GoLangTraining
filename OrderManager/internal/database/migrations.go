package database

import (
	"database/sql"
	"fmt"
	"log"
)

func RunMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")

	// First, check if the table exists
	var exists bool
	err := db.QueryRow(`
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_name = 'orders'
        );
    `).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error checking if table exists: %v", err)
	}

	if exists {
		log.Println("Orders table already exists, skipping migrations")
		return nil
	}

	// Run migrations within a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Create extensions and types
	_, err = tx.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		return fmt.Errorf("error creating uuid extension: %v", err)
	}

	// Create payment_status enum
	_, err = tx.Exec(`
        DO $$ 
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment_status') THEN
                CREATE TYPE payment_status AS ENUM ('paid', 'pending');
            END IF;
        END $$;
    `)
	if err != nil {
		return fmt.Errorf("error creating payment_status enum: %v", err)
	}

	// Create order_status enum
	_, err = tx.Exec(`
        DO $$ 
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
                CREATE TYPE order_status AS ENUM ('processing', 'pending', 'completed', 'cancelled');
            END IF;
        END $$;
    `)
	if err != nil {
		return fmt.Errorf("error creating order_status enum: %v", err)
	}

	// Create orders table
	_, err = tx.Exec(`
        CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inventory_id UUID NOT NULL,
    user_id UUID NOT NULL,
    payment_status payment_status NOT NULL DEFAULT 'pending',
    order_status order_status NOT NULL DEFAULT 'pending',
    price DECIMAL(10,2) NOT NULL,
    order_count INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
    `)
	if err != nil {
		return fmt.Errorf("error creating orders table: %v", err)
	}

	// Create indexes
	_, err = tx.Exec(`
        CREATE INDEX idx_user_id ON orders(user_id);
        CREATE INDEX idx_inventory_id ON orders(inventory_id);
    `)
	if err != nil {
		return fmt.Errorf("error creating indexes: %v", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
