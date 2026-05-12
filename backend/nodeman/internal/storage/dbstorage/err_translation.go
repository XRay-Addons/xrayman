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

	// if it is not our xerr (wrapped with stack previously)
	// wrap it with stack now
	// xerr.xerr - struct, error impl
	err = xerr.WrapWithStack(err)

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
	var pgErrCode string
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
		pgErrCode = pgErr.Code
	}

	// this is temporary sql error, we decided.
	// save as much details as possible by
	// extraction via %+v
	err = xerr.WrapWithf(errdefs.ErrTemporaryUnavailable, "%+v", err)
	// add pg error code if extracted
	if pgErrCode != "" {
		err = xerr.WrapWithf(err, "pg error code: %s", pgErrCode)
	}

	return err
}
