-- +goose Up
CREATE TABLE IF NOT EXISTS whatsapp_templates (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL,
    name          VARCHAR(255) NOT NULL,
    body          TEXT NOT NULL DEFAULT '',
    modified_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS phone_lists (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER NOT NULL,
    name          VARCHAR(255) NOT NULL,
    modified_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS phone_numbers (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    phone_list_id INTEGER NOT NULL,
    number       VARCHAR(32) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_phone_numbers_list ON phone_numbers(phone_list_id);

-- +goose Down
DROP TABLE IF EXISTS phone_numbers;
DROP TABLE IF EXISTS phone_lists;
DROP TABLE IF EXISTS whatsapp_templates;
