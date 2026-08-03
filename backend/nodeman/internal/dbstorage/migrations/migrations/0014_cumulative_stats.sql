-- +goose Up
-- +goose StatementBegin

DROP TABLE IF EXISTS daily_users_stats;
DROP TABLE IF EXISTS total_users_stats;

CREATE TABLE total_users_stats (
    stat_id BIGSERIAL PRIMARY KEY, 
    user_id    bigint    NOT NULL,
    download   bigint    NOT NULL,
    upload     bigint    NOT NULL
);

CREATE TABLE daily_users_stats (
    day        date      NOT NULL,
    node_id    bigint    NOT NULL,
    user_id    bigint    NOT NULL,
    download   bigint    NOT NULL,
    upload     bigint    NOT NULL,

    PRIMARY KEY (day, node_id, user_id)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS daily_users_stats;
DROP TABLE IF EXISTS total_users_stats;
-- +goose StatementEnd