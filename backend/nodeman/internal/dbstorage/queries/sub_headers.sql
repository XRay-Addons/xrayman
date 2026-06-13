-- name: NewSubHeader :one
	INSERT INTO sub_headers (
		header_key,
		header_value
	)
	VALUES ($1, $2)
	ON CONFLICT (header_key)
	DO UPDATE SET
		header_value = EXCLUDED.header_value,
		updated_at = now(),
		deleted_at = NULL
	RETURNING header_id;

-- name: DeleteSubHeader :exec
	UPDATE sub_headers
	SET deleted_at = now()
	WHERE header_id = $1
		AND deleted_at IS NULL;

-- name: ListSubHeaders :many
	SELECT
		header_id,
		header_key,
		header_value
	FROM sub_headers
	WHERE deleted_at IS NULL
	ORDER BY header_id ASC;