-- +goose Up
-- +goose StatementBegin

ALTER TABLE dynamic_config
    RENAME TO settings;

ALTER TABLE settings
    RENAME COLUMN config TO settings;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE settings
    RENAME COLUMN settings TO config;

ALTER TABLE settings
    RENAME TO dynamic_config;

-- +goose StatementEnd