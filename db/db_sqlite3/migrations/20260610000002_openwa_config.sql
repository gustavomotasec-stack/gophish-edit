-- +goose Up
CREATE TABLE IF NOT EXISTS openwa_configs (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER NOT NULL UNIQUE,
    api_url     VARCHAR(255) NOT NULL DEFAULT 'http://localhost:3000',
    api_key     VARCHAR(255) NOT NULL DEFAULT '',
    min_delay   INTEGER NOT NULL DEFAULT 3,
    max_delay   INTEGER NOT NULL DEFAULT 8
);

-- +goose Down
DROP TABLE IF EXISTS openwa_configs;
