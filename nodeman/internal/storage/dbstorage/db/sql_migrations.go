package db

import (
	"context"
	"database/sql"
	"embed"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func migrate(ctx context.Context, db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return errdefs.WrapWithStack(err)
	}
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return errdefs.WrapWithStack(err)
	}
	return nil
}
