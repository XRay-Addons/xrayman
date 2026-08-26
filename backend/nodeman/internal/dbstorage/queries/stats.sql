-- name: SetLocalTxFastMode :exec
SET LOCAL synchronous_commit = OFF;

-- name: UpdateStats :exec
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
INSERT INTO nodes_stats (
    node_id,
    upload,
    download,
    open_connections,
    cpu_load,
    ram_load,
    mem_load
)
SELECT 
    sqlc.arg(node_id)::bigint,
    COALESCE(SUM(upload), 0)::bigint, 
    COALESCE(SUM(download), 0)::bigint,
    sqlc.arg(open_connections)::integer,
    sqlc.arg(cpu_load)::real,
    sqlc.arg(ram_load)::real,
    sqlc.arg(mem_load)::real
FROM input_data
ON CONFLICT (node_id) DO UPDATE
SET
    upload           = nodes_stats.upload   + EXCLUDED.upload,
    download         = nodes_stats.download + EXCLUDED.download,
    open_connections = EXCLUDED.open_connections,
    cpu_load         = EXCLUDED.cpu_load,
    ram_load         = EXCLUDED.ram_load,
    mem_load         = EXCLUDED.mem_load;

-- name: RefreshDailyStats :exec
WITH
    -- update users daily traffic stats
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
    -- update nodes daily traffic stats
    snapshot_nodes AS (
        INSERT INTO daily_nodes_traffic (day, node_id, upload, download)
    SELECT
        sqlc.arg(day)::date AS day,
        t.node_id,
        t.upload,
        t.download
    FROM nodes_stats t
    ON CONFLICT (day, node_id) DO UPDATE
    SET
        upload   = GREATEST(daily_nodes_traffic.upload, EXCLUDED.upload),
        download = GREATEST(daily_nodes_traffic.download, EXCLUDED.download)
    )
-- do nothing stub
SELECT 1;