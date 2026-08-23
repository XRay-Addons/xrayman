
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
WHERE n.deleted_at IS NULL
  AND u.revision > n.revision;


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
    AND n.deleted_at IS NULL;
