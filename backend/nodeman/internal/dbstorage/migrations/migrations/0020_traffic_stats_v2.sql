-- +goose Up
-- +goose StatementBegin

DROP TABLE IF EXISTS total_users_traffic;
DROP TABLE IF EXISTS daily_users_traffic;
DROP TABLE IF EXISTS total_nodes_traffic;
DROP TABLE IF EXISTS daily_nodes_traffic;
DROP INDEX IF EXISTS daily_users_traffic_index;
DROP INDEX IF EXISTS daily_nodes_traffic_index;

CREATE TABLE total_users_traffic (
    user_id    bigint    NOT NULL,
    download   bigint    NOT NULL,
    upload     bigint    NOT NULL,

    PRIMARY KEY (user_id)
);

CREATE TABLE daily_users_traffic (
    day        date      NOT NULL,
    user_id    bigint    NOT NULL,
    download   bigint    NOT NULL,
    upload     bigint    NOT NULL,

    PRIMARY KEY (day, user_id)
);

CREATE INDEX daily_users_traffic_index ON daily_users_traffic (user_id, day DESC);

CREATE TABLE total_nodes_traffic (
    node_id    bigint    NOT NULL,
    download   bigint    NOT NULL,
    upload     bigint    NOT NULL,

    PRIMARY KEY (node_id)
);

CREATE TABLE daily_nodes_traffic (
    day        date      NOT NULL,
    node_id    bigint    NOT NULL,
    download   bigint    NOT NULL,
    upload     bigint    NOT NULL,

    PRIMARY KEY (day, node_id)
);

CREATE INDEX daily_nodes_traffic_index ON daily_nodes_traffic (node_id, day DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS total_users_traffic;
DROP TABLE IF EXISTS daily_users_traffic;
DROP TABLE IF EXISTS total_nodes_traffic;
DROP TABLE IF EXISTS daily_nodes_traffic;
DROP INDEX IF EXISTS daily_users_traffic_index;
DROP INDEX IF EXISTS daily_nodes_traffic_index;
-- +goose StatementEnd