-- +goose Up
-- +goose StatementBegin

CREATE INDEX total_users_traffic_idx ON total_users_traffic (user_id);

CREATE INDEX daily_users_traffic_index ON daily_users_traffic (user_id, day DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS daily_users_traffic_index;
-- +goose StatementEnd