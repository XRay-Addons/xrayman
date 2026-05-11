package dbstorage

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
)

func (uow *uowctx) NewSubHeader(ctx context.Context, header *models.Header) error {
	query := queryReplacer.Replace(`
		INSERT INTO {sub_headers} (
			{header_key},
			{header_val},
		) VALUES ($1, $2)
		RETURNING {header_id}
	`)

	err := uow.tx.QueryRowContext(ctx, query,
		header.Key,
		header.Value,
	).Scan(&header.ID)

	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func (uow *uowctx) DeleteSubHeader(ctx context.Context, id models.HeaderID) error {
	query := queryReplacer.Replace(`
		UPDATE {sub_headers}
		SET {deleted_at} = now()
		WHERE {header_id} = $1
		  AND {deleted_at} IS NULL
	`)

	_, err := uow.tx.ExecContext(ctx, query, id)
	if err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func (uow *uowctx) ListSubHeaders(ctx context.Context) ([]models.Header, error) {
	query := queryReplacer.Replace(`
		SELECT
			{header_id},
			{header_key},
			{header_val},
		FROM {sub_headers}
		WHERE {deleted_at} IS NULL
		ORDER BY {header_id} ASC
	`)

	rows, err := uow.tx.QueryContext(ctx, query)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	defer rows.Close()

	var headers []models.Header
	for rows.Next() {
		var header models.Header
		err := rows.Scan(
			&header.ID,
			&header.Key,
			&header.Value,
		)
		if err != nil {
			return nil, xerr.WrapWithStack(err)
		}
		headers = append(headers, header)
	}

	if err := rows.Err(); err != nil {
		return nil, xerr.WrapWithStack(err)
	}

	return headers, nil
}
