-- name: SetNodeRev :exec
UPDATE nodes
SET
    revision = $1,
    updated_at = now()
WHERE node_id = $2;

-- name: FindPendingSyncs :many
SELECT
    u.user_id,
    u.user_name,
    u.display_name,
    u.vless_uuid,
    u.user_target_status,
    u.revision
FROM users AS u
JOIN nodes AS n
    ON n.node_id = $1
WHERE u.revision > n.revision
ORDER BY u.user_id ASC; 

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
INNER JOIN users u
    ON u.user_id = $1
WHERE n.revision >= u.revision
    AND n.node_current_status = sqlc.arg(node_status_running)
    AND n.node_target_status = sqlc.arg(node_status_running)
    AND n.deleted_at IS NULL
ORDER BY n.node_id ASC; 
