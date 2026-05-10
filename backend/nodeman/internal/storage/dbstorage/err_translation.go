package dbstorage

import (
	"context"

	"errors"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	errcode "github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func translatePgErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	if errors.Is(err, context.Canceled) {
		return err
	}

	var pgParseErr *pgconn.ParseConfigError
	if errors.As(err, &pgParseErr) {
		// invalid config, unreteyable
		return err
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// whitelist of unretryable pg errors
		// for incorrect requests
		switch pgErr.Code {
		case errcode.SyntaxError,
			errcode.UndefinedTable,
			errcode.InvalidPassword,
			errcode.InvalidCatalogName:
			{
				return err
			}
		}
	}

	return xerr.Wrap(errdefs.ErrTemporaryUnavailable,
		xerr.WithStack(), xerr.With(err.Error()))
}
