CREATE TABLE IF NOT EXISTS banners (
    id BIGSERIAL PRIMARY KEY,
    is_active BOOLEAN NOT NULL,
    feature BIGINT NOT NULL,
    content BYTEA,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL
);