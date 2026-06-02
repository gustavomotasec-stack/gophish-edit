-- +goose Up
ALTER TABLE campaigns ADD COLUMN campaign_type VARCHAR(255) NOT NULL DEFAULT 'email';

-- +goose Down
-- SQLite does not support DROP COLUMN.
