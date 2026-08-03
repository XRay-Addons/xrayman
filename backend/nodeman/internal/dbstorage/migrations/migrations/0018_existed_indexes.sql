-- +goose Up
-- +goose StatementBegin
CREATE INDEX nodes_existed_index ON nodes (node_id) WHERE deleted_at IS NULL;
CREATE INDEX user_existed_index ON users (user_id) WHERE deleted_at IS NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS nodes_existed_index;
DROP INDEX IF EXISTS user_existed_index;
-- +goose StatementEnd