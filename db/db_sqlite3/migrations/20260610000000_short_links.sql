-- +goose Up
CREATE TABLE IF NOT EXISTS short_links (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    code        VARCHAR(16) NOT NULL UNIQUE,
    original    TEXT NOT NULL,
    campaign_id INTEGER NOT NULL DEFAULT 0,
    rid         VARCHAR(32) NOT NULL DEFAULT '',
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_short_links_code ON short_links(code);

-- +goose Down
DROP TABLE IF EXISTS short_links;
