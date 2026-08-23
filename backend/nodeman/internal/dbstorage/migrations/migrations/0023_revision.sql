-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE global_revision_seq;
ALTER TABLE users ADD COLUMN revision BIGINT NOT NULL DEFAULT 1;
ALTER TABLE nodes ADD COLUMN revision BIGINT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE nodes DROP COLUMN revision;
ALTER TABLE users DROP COLUMN revision;
DROP SEQUENCE global_revision_seq;
-- +goose StatementEnd

