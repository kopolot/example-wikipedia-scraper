CREATE TYPE order_status AS ENUM ('pending', 'paid', 'completed', 'refunded', 'cancelled' , 'failed');
CREATE TYPE order_item_type AS ENUM ('subscription');

CREATE TABLE "orders" (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    client_ip VARCHAR(255) NOT NULL,
    status order_status NOT NULL,
    invoice_number VARCHAR(255) UNIQUE,
    currency VARCHAR(10) NOT NULL DEFAULT 'EUR',
    total_amount NUMERIC(10, 2) NOT NULL,
    payment_id VARCHAR(255),
    payment_method VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_order_user ON "orders" (user_id);
CREATE INDEX idx_order_status ON "orders" (status);
CREATE INDEX idx_order_payment_id ON "orders" (payment_id);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES "orders"(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    product_id INTEGER NOT NULL REFERENCES products(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    price NUMERIC(10, 2) NOT NULL,
    quantity INTEGER NOT NULL,
    item_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_order_item_order ON order_items (order_id);