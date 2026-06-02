-- +goose Up
ALTER TABLE results ADD COLUMN source VARCHAR(255) NOT NULL DEFAULT 'email';

-- +goose Down
-- SQLite does not support DROP COLUMN, so this migration cannot be reversed.
