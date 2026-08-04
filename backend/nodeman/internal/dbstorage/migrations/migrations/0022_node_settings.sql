-- +goose Up
-- +goose StatementBegin
ALTER TABLE nodes
ADD COLUMN version TEXT NOT NULL DEFAULT 'v0.0.0';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE nodes
DROP COLUMN version;
-- +goose StatementEnd