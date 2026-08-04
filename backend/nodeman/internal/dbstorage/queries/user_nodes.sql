
-- name: FindPendingSyncs :many
SELECT
    u.user_id,
    u.user_name,
    u.display_name,
    u.vless_uuid,
    u.user_target_status,
    COALESCE(
        s.user_current_status,
        sqlc.arg(default_user_status)::smallint
    ) AS user_current_status
FROM users u
LEFT JOIN syncs s
    ON s.user_id = u.user_id
   AND s.node_id = $1
WHERE
    COALESCE(
        s.user_current_status,
        sqlc.arg(default_user_status)::smallint
    ) IS DISTINCT FROM u.user_target_status;

-- name: DeleteNodeUsers :exec
DELETE FROM syncs
WHERE node_id = $1;

-- name: InsertNodeUsers :exec
INSERT INTO syncs (user_id, node_id, user_current_status)
SELECT
    t.user_id,
    sqlc.arg(node_id)::bigint,
    t.user_current_status
FROM ROWS FROM (
    unnest(sqlc.arg(user_id)::bigint[]),
    unnest(sqlc.arg(user_current_status)::smallint[])
) AS t(user_id, user_current_status)
ON CONFLICT (user_id, node_id)
DO UPDATE SET user_current_status = EXCLUDED.user_current_status;

-- name: GetUserNodes :many
SELECT
    n.node_id,
    n.client_cfg_template,
    n.version,
    n.node_endpoint,
    n.node_access_key,
    n.node_current_status,
    n.node_target_status
FROM nodes n
INNER JOIN syncs s
    ON s.node_id = n.node_id
WHERE s.user_id = $1
    AND s.user_current_status = sqlc.arg(user_status_enabled)::smallint
    AND n.node_target_status = sqlc.arg(node_status_running)::smallint
    AND n.node_current_status = sqlc.arg(node_status_running)::smallint
    AND n.deleted_at IS NULL;
