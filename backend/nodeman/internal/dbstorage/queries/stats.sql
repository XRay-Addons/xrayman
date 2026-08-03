-- name: SetLocalTxFastMode :exec
SET LOCAL synchronous_commit = OFF;

-- name: UpdateTotalStats :exec
WITH 
-- 1. input -> flat table
input_data AS (
    SELECT 
        t.user_id,
        t.upload,
        t.download
    FROM ROWS FROM (
        unnest(sqlc.arg(user_id)::bigint[]),
        unnest(sqlc.arg(upload)::bigint[]),
        unnest(sqlc.arg(download)::bigint[])
    ) AS t(user_id, upload, download)
),
-- 2. update users traffic
update_users AS (
    INSERT INTO total_users_traffic (user_id, upload, download)
    SELECT user_id, upload, download FROM input_data
    ON CONFLICT (user_id) DO UPDATE
    SET
        upload   = total_users_traffic.upload   + EXCLUDED.upload,
        download = total_users_traffic.download + EXCLUDED.download
    RETURNING 1 -- sqlc require to return something
)
-- 3. update nodes stats
INSERT INTO total_nodes_traffic (node_id, upload, download)
SELECT 
    sqlc.arg(node_id)::bigint,
    COALESCE(SUM(upload), 0)::bigint, 
    COALESCE(SUM(download), 0)::bigint
FROM input_data
ON CONFLICT (node_id) DO UPDATE
SET
    upload   = total_nodes_traffic.upload   + EXCLUDED.upload,
    download = total_nodes_traffic.download + EXCLUDED.download;

-- name: UpdateDailyStats :exec
WITH
    -- update users daily stats
    snapshot_users AS (
        INSERT INTO daily_users_traffic (day, user_id, upload, download)
        SELECT
            sqlc.arg(day)::date AS day,
            t.user_id,
            t.upload,
            t.download
        FROM total_users_traffic t
        ON CONFLICT (day, user_id) DO UPDATE
        SET
            upload   = GREATEST(daily_users_traffic.upload, EXCLUDED.upload),
            download = GREATEST(daily_users_traffic.download, EXCLUDED.download)
        RETURNING 1
    ),
    -- update nodes daily stats
    snapshot_nodes AS (
        INSERT INTO daily_nodes_traffic (day, node_id, upload, download)
    SELECT
        sqlc.arg(day)::date AS day,
        t.node_id,
        t.upload,
        t.download
    FROM total_nodes_traffic t
    ON CONFLICT (day, node_id) DO UPDATE
    SET
        upload   = GREATEST(daily_nodes_traffic.upload, EXCLUDED.upload),
        download = GREATEST(daily_nodes_traffic.download, EXCLUDED.download)
    )
-- do nothing stub
SELECT 1;