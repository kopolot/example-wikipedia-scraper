CREATE TABLE subscription_levels (
    id SERIAL PRIMARY KEY,
    "level" integer NOT NULL UNIQUE,
    "limit" integer NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE subscription_level_products (
    subscription_level_id integer NOT NULL,
    product_id integer NOT NULL,
    PRIMARY KEY (subscription_level_id, product_id),
    FOREIGN KEY (subscription_level_id) REFERENCES subscription_levels(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);