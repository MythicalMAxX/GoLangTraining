-- Up migration
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inventory_id UUID NOT NULL,
    user_id UUID NOT NULL,
    payment_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    order_status VARCHAR(50) NOT NULL DEFAULT 'pending',
    price DECIMAL(10,2) NOT NULL,
    order_count INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on user_id for faster lookups
CREATE INDEX idx_orders_user_id ON orders(user_id);

-- Create an index on inventory_id for faster lookups
CREATE INDEX idx_orders_inventory_id ON orders(inventory_id);

-- Down migration
DROP TABLE IF EXISTS orders;