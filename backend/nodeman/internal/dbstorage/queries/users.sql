-- name: NewUser :one
INSERT INTO users (
    display_name,
    user_name,
    vless_uuid,
    user_target_status,
    revision
) VALUES (
    $1,
    $2,
    $3,
    $4,
    nextval('global_revision_seq')
) RETURNING user_id;

-- name: GetUserView :one
SELECT
    u.user_id,
    u.display_name,
    u.user_name,
    u.vless_uuid,
    u.user_target_status,

    COALESCE(total_stats.upload, 0)   AS upload_total,
    COALESCE(total_stats.download, 0) AS download_total,

    (COALESCE(total_stats.upload, 0)
      - COALESCE(daily_stats.upload, 0))::bigint AS upload_last_days,

    (COALESCE(total_stats.download, 0)
      - COALESCE(daily_stats.download, 0))::bigint AS download_last_days

FROM users u

LEFT JOIN (
    SELECT
        user_id,
        SUM(upload)   AS upload,
        SUM(download) AS download
    FROM total_users_traffic
    WHERE user_id = sqlc.arg(user_id)::bigint
    GROUP BY user_id
) total_stats ON total_stats.user_id = u.user_id

LEFT JOIN (
    SELECT
        upload,
        download
    FROM daily_users_traffic
    WHERE user_id = sqlc.arg(user_id)::bigint
      AND day < sqlc.arg(from_day)::date
    ORDER BY day DESC
    LIMIT 1
) daily_stats ON TRUE

WHERE u.deleted_at IS NULL
  AND u.user_id = sqlc.arg(user_id)::bigint
  AND u.user_name = sqlc.arg(user_name)::text;


-- name: ListUsers :many
SELECT
    u.user_id,
    u.display_name,
    u.user_name,
    u.vless_uuid,
    u.user_target_status
FROM users u
WHERE deleted_at IS NULL
ORDER BY u.user_id ASC;

-- name: ListUserViews :many
SELECT
    u.user_id,
    u.display_name,
    u.user_name,
    u.vless_uuid,
    u.user_target_status,

    COALESCE(total_stats.upload, 0)   AS upload_total,
    COALESCE(total_stats.download, 0) AS download_total,

    (COALESCE(total_stats.upload, 0)
      - COALESCE(daily_stats.upload, 0))::bigint AS upload_last_days,

    (COALESCE(total_stats.download, 0)
      -COALESCE(daily_stats.download, 0))::bigint AS download_last_days

FROM users u

LEFT JOIN (
    SELECT
        user_id,
        SUM(upload)   AS upload,
        SUM(download) AS download
    FROM total_users_traffic
    GROUP BY user_id
) total_stats ON total_stats.user_id = u.user_id

LEFT JOIN (
    SELECT DISTINCT ON (user_id)
        user_id,
        upload,
        download,
        day
    FROM daily_users_traffic
    WHERE day < sqlc.arg(from_day)::date
    ORDER BY user_id, day DESC
) daily_stats ON daily_stats.user_id = u.user_id

WHERE u.deleted_at IS NULL
ORDER BY u.user_id ASC;

-- name: SetTargetUserStatus :exec
UPDATE users
SET
    user_target_status = $1,
    revision = nextval('global_revision_seq'),
    updated_at = now()
WHERE user_id = $2
    AND deleted_at IS NULL
    AND user_target_status IS DISTINCT FROM $1;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE user_id = $1
    AND deleted_at IS NULL;