CREATE TABLE "users" (
    id SERIAL PRIMARY KEY,
    last_login_at TIMESTAMPTZ,
    password_reset_token_expires_at TIMESTAMPTZ DEFAULT NULL,
    last_logout_at TIMESTAMPTZ DEFAULT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL DEFAULT 'user',
    username VARCHAR(255) NOT NULL UNIQUE,
    email_verification_token VARCHAR(255),
    password_reset_token VARCHAR(255),
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);