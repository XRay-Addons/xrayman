-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS nodes (
    node_id BIGSERIAL PRIMARY KEY,
    client_cfg_template TEXT NOT NULL,
    node_endpoint TEXT NOT NULL,
    node_access_key BYTEA NOT NULL CHECK (octet_length(node_access_key) = 64),
    node_current_status SMALLINT NOT NULL DEFAULT 0,
    node_target_status SMALLINT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS nodes;
-- +goose StatementEnd