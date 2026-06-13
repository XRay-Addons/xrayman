-- +goose Up
-- +goose StatementBegin

DELETE FROM sub_headers a
WHERE a.deleted_at IS NOT NULL;

DELETE FROM sub_headers a
USING sub_headers b
WHERE a.ctid < b.ctid
  AND a.header_key = b.header_key;

ALTER TABLE sub_headers
ADD CONSTRAINT sub_headers_header_key_unique UNIQUE (header_key);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sub_headers
DROP CONSTRAINT IF EXISTS sub_headers_header_key_unique;
-- +goose StatementEnd