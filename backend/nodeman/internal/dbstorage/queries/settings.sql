-- name: EnsureSettings :exec
INSERT INTO settings (id, settings, updated_at)
VALUES (true, sqlc.arg(cfg)::jsonb, now())
ON CONFLICT (id)
DO NOTHING;

-- name: GetSettings :one
SELECT settings
FROM settings
WHERE id = true;

-- name: SetSettings :exec
INSERT INTO settings (id, settings, updated_at)
VALUES (true, sqlc.arg(cfg)::jsonb, now())
ON CONFLICT (id)
DO UPDATE SET
    settings = EXCLUDED.settings,
    updated_at = now();

