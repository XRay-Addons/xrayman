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

-- name: ListNodeViews :many
SELECT
    n.node_id,
    n.client_cfg_template,
    n.version,
    n.node_endpoint,
    n.node_access_key,
    n.node_current_status,
    n.node_target_status,

    COALESCE(ts.upload, 0)   AS upload_total,
    COALESCE(ts.download, 0) AS download_total,

    (COALESCE(ts.upload, 0) - COALESCE(ds.upload, 0))::bigint AS upload_last_days,
    (COALESCE(ts.download, 0) - COALESCE(ds.download, 0))::bigint AS download_last_days,

    COALESCE(ts.open_connections, 0)  AS open_connections,
    COALESCE(ts.cpu_load, 0)          AS cpu_load,
    COALESCE(ts.ram_load, 0)          AS ram_load,
    COALESCE(ts.mem_load, 0)          AS mem_load
FROM nodes n

LEFT JOIN (
    SELECT
        node_id,
        upload,
        download,
        open_connections,
        cpu_load,
        ram_load,
        mem_load
    FROM nodes_stats
    GROUP BY node_id
) ts ON ts.node_id = n.node_id

LEFT JOIN (
    SELECT DISTINCT ON (node_id)
        node_id,
        upload,
        download,
        day
    FROM daily_nodes_traffic
    WHERE day < sqlc.arg(from_day)::date
    ORDER BY node_id, day DESC
) ds ON ds.node_id = n.node_id

WHERE deleted_at IS NULL
ORDER BY n.node_id ASC;

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
