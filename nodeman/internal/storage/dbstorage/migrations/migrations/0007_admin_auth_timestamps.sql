-- +goose Up
-- +goose StatementBegin
ALTER TABLE admin_auth
	ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	ADD COLUMN deleted_at TIMESTAMPTZ;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE admin_auth
	DROP COLUMN deleted_at,
	DROP COLUMN updated_at,
	DROP COLUMN created_at;
-- +goose StatementEnd