CREATE TABLE IF NOT EXISTS banners (
    id SERIAL PRIMARY KEY,
    is_active BOOLEAN,
    feature INT,
    content BYTEA,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);