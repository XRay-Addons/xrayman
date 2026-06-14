-- +goose Up
-- +goose StatementBegin

CREATE TABLE dynamic_config (
    id          boolean PRIMARY KEY DEFAULT true,
    config      jsonb NOT NULL,
    updated_at  timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT single_row CHECK (id = true)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dynamic_config;
-- +goose StatementEnd