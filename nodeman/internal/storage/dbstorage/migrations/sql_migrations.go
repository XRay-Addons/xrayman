package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type option func(options *options)

type options struct {
	retry bool
	log   *zap.Logger
}

/*func WithRetry() option {
	return func(o *options) {
		o.retry = true
	}
}

func WithLogger(log *zap.Logger) option {
	return func(o *options) {
		if log != nil {
			o.log = log
		}
	}
}*/

func ApplyMigrations(ctx context.Context, db *sql.DB, log *zap.Logger) error {
	if db == nil {
		return errdefs.NewNilArg("db")
	}

	goose.SetBaseFS(embedMigrations)
	goose.SetLogger(gooseLogger(log))
	defer func() {
		goose.SetBaseFS(nil)
		goose.SetLogger(nil)
	}()

	if err := goose.SetDialect("postgres"); err != nil {
		return errdefs.WrapWithStack(err)
	}

	//if !retry {
	// migrate without retries
	return migrate(ctx, db)
	//}

	/*// retry with policy till success or cancel
	const inintalRetry = 100 * time.Millisecond
	const maxRetry = 2 * time.Second
	b := goretry.WithMaxDuration(maxRetry,
		goretry.NewFibonacci(inintalRetry))

	return goretry.Do(ctx, b, func(ctx context.Context) error {
		if err := migrate(ctx, db); err != nil {
			log.Warn("migration attempt", zap.Error(err))
			return err
		}
		log.Warn("migration successed")
		return nil
	})*/
}

// zap.logger to goose.logger adapter
type gl struct {
	l *zap.Logger
}

func (g *gl) Fatalf(format string, v ...interface{}) {
	g.l.Fatal(fmt.Sprintf(format, v...))
}

func (g *gl) Printf(format string, v ...interface{}) {
	g.l.Info(fmt.Sprintf(format, v...))
}

var _ goose.Logger = (*gl)(nil)

func gooseLogger(l *zap.Logger) goose.Logger {
	return &gl{l: l}
}

func migrate(ctx context.Context, db *sql.DB) error {
	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return errdefs.WrapWithStack(err)
	}
	return nil
}
