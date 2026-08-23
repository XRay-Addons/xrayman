-- name: NewNode :one
INSERT INTO nodes (
    client_cfg_template,
    version,
    node_endpoint,
    node_access_key,
    node_current_status,
    node_target_status
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING node_id;

-- name: GetNode :one
SELECT
    node_id,
    client_cfg_template,
    version,
    node_endpoint,
    node_access_key,
    node_current_status,
    node_target_status
FROM nodes
WHERE node_id = $1
    AND deleted_at IS NULL;

-- name: ListNodes :many
SELECT
    node_id,
    client_cfg_template,
    version,
    node_endpoint,
    node_access_key,
    node_current_status,
    node_target_status
FROM nodes
WHERE deleted_at IS NULL
ORDER BY node_id ASC;

-- name: SetTargetNodeStatus :exec
UPDATE nodes
SET
    node_target_status = $1,
    updated_at = now()
WHERE node_id = $2
    AND deleted_at IS NULL;

-- name: SetCurrentNodeStatus :exec
UPDATE nodes
SET
    node_current_status = $1,
    updated_at = now()
WHERE node_id = $2
    AND deleted_at IS NULL;

-- name: SetNodeSettings :exec
UPDATE nodes
SET
    client_cfg_template = $1,
    version = $2,
    updated_at = now()
WHERE node_id = $3
    AND deleted_at IS NULL;

-- name: DeleteNode :exec
UPDATE nodes
SET deleted_at = now()
WHERE node_id = $1
    AND deleted_at IS NULL;

-- name: SetNodeRev :exec
UPDATE nodes
SET
    revision = $1,
    updated_at = now()
WHERE node_id = $2;

-- name: GetNodeRev :one
SELECT
   revision
FROM nodes
WHERE node_id = $1
    AND deleted_at IS NULL;