-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS syncs (
    user_id BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    node_id BIGINT NOT NULL REFERENCES nodes(node_id) ON DELETE CASCADE,
    user_current_status SMALLINT NOT NULL DEFAULT 0,
    PRIMARY KEY (user_id, node_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS syncs;
-- +goose StatementEnd
