-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS admin_auth (
    admin_id BIGSERIAL PRIMARY KEY,
    password_hash BYTEA NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS admin_auth;
-- +goose StatementEnd
