-- name: GetPassword :one
SELECT
    admin_id,
    password_hash
FROM admin_auth
WHERE admin_id = $1
    AND deleted_at IS NULL;

-- name: SetPassword :exec
INSERT INTO admin_auth (
    admin_id,
    password_hash,
    updated_at
) VALUES ($1, $2, now())
ON CONFLICT (admin_id)
DO UPDATE
SET
    password_hash = EXCLUDED.password_hash,
    updated_at = now();