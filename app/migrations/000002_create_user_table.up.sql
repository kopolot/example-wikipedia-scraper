CREATE TABLE "users" (
    id SERIAL PRIMARY KEY,
    last_login_at TIMESTAMPTZ,
    password_reset_token_expires_at TIMESTAMPTZ DEFAULT NULL,
    last_logout_at TIMESTAMPTZ DEFAULT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    username TEXT NOT NULL UNIQUE,
    email_verification_token TEXT,
    password_reset_token TEXT,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);