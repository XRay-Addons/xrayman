-- +goose Up
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN visible_name TO display_name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users RENAME COLUMN display_name TO visible_name;
-- +goose StatementEnd