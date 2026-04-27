-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
	ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	ADD COLUMN deleted_at TIMESTAMPTZ;

ALTER TABLE nodes
	ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	ADD COLUMN deleted_at TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE nodes
	DROP COLUMN deleted_at,
	DROP COLUMN updated_at,
	DROP COLUMN created_at;

ALTER TABLE users
	DROP COLUMN deleted_at,
	DROP COLUMN updated_at,
	DROP COLUMN created_at;
-- +goose StatementEnd
