CREATE TABLE "pages" (
    id SERIAL PRIMARY KEY,
    site_name TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    content TEXT,
    title TEXT NOT NULL,
    text_field_1 TEXT,
    text_field_2 TEXT,
    text_field_3 TEXT,
    hash_key TEXT NOT NULL UNIQUE,
    external_id TEXT NOT NULL,
    notified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    CONSTRAINT idx_page_site_external UNIQUE (site_name, url, external_id),
    CONSTRAINT idx_page_site_hash UNIQUE (site_name, hash_key)
);