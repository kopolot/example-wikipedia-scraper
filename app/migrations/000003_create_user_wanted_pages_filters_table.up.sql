CREATE TABLE user_wanted_pages_filters (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL DEFAULT '',
    site_names TEXT[] NULL,
    keywords TEXT[] NULL,
    title_contains VARCHAR(512) NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_user_wanted_pages_filters_user_id ON user_wanted_pages_filters(user_id);
CREATE INDEX idx_user_wanted_pages_filters_site_names ON user_wanted_pages_filters USING GIN (site_names);
CREATE INDEX idx_user_wanted_pages_filters_keywords ON user_wanted_pages_filters USING GIN (keywords);
