-- name: EnsureDynamicConfig :exec
INSERT INTO dynamic_config (id, config, updated_at)
VALUES (true, sqlc.arg(cfg)::jsonb, now())
ON CONFLICT (id)
DO NOTHING;

-- name: GetDynamicConfig :one
SELECT config
FROM dynamic_config
WHERE id = true;

-- name: SetDynamicConfig :exec
INSERT INTO dynamic_config (id, config, updated_at)
VALUES (true, sqlc.arg(cfg)::jsonb, now())
ON CONFLICT (id)
DO UPDATE SET
    config = EXCLUDED.config,
    updated_at = now();

