-- +goose Up
-- +goose StatementBegin

CREATE TABLE total_users_stats (
    user_id     BIGINT NOT NULL,
    node_id     BIGINT NOT NULL,
    upload      BIGINT NOT NULL DEFAULT 0,
    download    BIGINT NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, node_id)
);

CREATE TABLE daily_users_stats (
    user_id   BIGINT NOT NULL,
    node_id   BIGINT NOT NULL,
    day       DATE   NOT NULL,
    upload    BIGINT NOT NULL DEFAULT 0,
    download  BIGINT NOT NULL DEFAULT 0,

    PRIMARY KEY (user_id, node_id, day)
);

CREATE INDEX idx_daily_user_day
    ON daily_users_stats (user_id, day);
    
-- +goose StatementEnd

DROP TABLE IF EXISTS daily_users_stats;
DROP TABLE IF EXISTS total_users_stats;

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd