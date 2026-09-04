-- name: GetMetrics :many
SELECT
    n.node_id,
    COALESCE(ns.download, 0)          AS download,
    COALESCE(ns.upload, 0)            AS upload,
    COALESCE(ns.open_connections, 0)  AS open_connections,
    COALESCE(ns.cpu_load, 0)          AS cpu_load,
    COALESCE(ns.ram_load, 0)          AS ram_load,
    COALESCE(ns.mem_load, 0)          AS mem_load,
    n.node_endpoint
FROM nodes n
LEFT JOIN nodes_stats ns ON ns.node_id = n.node_id
WHERE n.node_current_status = sqlc.arg(node_status_running)
    AND n.node_target_status = sqlc.arg(node_status_running)
    AND n.deleted_at IS NULL
ORDER BY n.node_id ASC; 