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
-- 1. settings -> now - (recent days count)
WITH from_day AS (
    SELECT
        (CURRENT_DATE- COALESCE((settings->>'RecentDays')::int, 0))::date AS value
    FROM settings
)
-- 2. select view from users
SELECT
    u.user_id,
    u.display_name,
    u.user_name,
    u.vless_uuid,
    u.user_target_status,

    COALESCE(ts.upload, 0)   AS upload_total,
    COALESCE(ts.download, 0) AS download_total,

    (COALESCE(ts.upload, 0) - COALESCE(ds.upload, 0))::bigint AS upload_recent_days,
    (COALESCE(ts.download, 0) - COALESCE(ds.download, 0))::bigint AS download_recent_days
FROM users u
-- 3. merged with stats
LEFT JOIN (
    SELECT
        user_id,
        upload,
        download
    FROM total_users_traffic
    WHERE user_id = sqlc.arg(user_id)::bigint
    GROUP BY user_id
) ts ON ts.user_id = u.user_id
-- 4. and mention recent days interval
LEFT JOIN (
    SELECT
        upload,
        download
    FROM daily_users_traffic
    CROSS JOIN from_day
    WHERE user_id = sqlc.arg(user_id)::bigint
      AND daily_users_traffic.day < from_day.value
    ORDER BY day DESC
    LIMIT 1
) ds ON TRUE

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
-- 1. settings -> now - (recent days count)
WITH from_day AS (
    SELECT
        (CURRENT_DATE- COALESCE((settings->>'RecentDays')::int, 0))::date AS value
    FROM settings
)
-- 2. select views from users
SELECT
    u.user_id,
    u.display_name,
    u.user_name,
    u.vless_uuid,
    u.user_target_status,

    COALESCE(ts.upload, 0)   AS upload_total,
    COALESCE(ts.download, 0) AS download_total,

    (COALESCE(ts.upload, 0) - COALESCE(ds.upload, 0))::bigint AS upload_recent_days,
    (COALESCE(ts.download, 0) - COALESCE(ds.download, 0))::bigint AS download_recent_days
FROM users u
-- 3. merged with stats
LEFT JOIN (
    SELECT
        user_id,
        upload,
        download
    FROM total_users_traffic
) ts ON ts.user_id = u.user_id
-- 4. and mention recent days interval
LEFT JOIN (
    SELECT DISTINCT ON (user_id)
        user_id,
        upload,
        download
    FROM daily_users_traffic
    CROSS JOIN from_day
    WHERE daily_users_traffic.day < from_day.value
    ORDER BY user_id, day DESC
) ds ON ds.user_id = u.user_id

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