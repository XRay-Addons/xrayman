-- name: Lock :exec
SELECT pg_advisory_xact_lock(sqlc.arg(lock_id)::bigint);