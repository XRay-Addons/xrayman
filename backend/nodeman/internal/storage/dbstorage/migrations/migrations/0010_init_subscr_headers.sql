-- +goose Up
-- +goose StatementBegin
CREATE TABLE sub_headers (
    header_id      BIGSERIAL PRIMARY KEY,
    header_key     TEXT NOT NULL,
    header_value   TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sub_headers;
-- +goose StatementEnd
