package dbstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	errcode "github.com/jackc/pgerrcode"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func translatePgErr(err error) error {
	if err == nil {
		return nil
	}

	// wrap with stack everything
	err = xerr.WrapWithStack(err)

	// context error
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	// connection string error
	var pgParseErr *pgconn.ParseConfigError
	if errors.As(err, &pgParseErr) {
		return err
	}

	// not found error
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
		return xerr.WrapWithType(err, errdefs.ErrNotFound)
	}

	// postgress syntax errors
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && isPgSyntaxError(pgErr.Code) {
		return err
	}

	// all other errors mark as temporary and retriable because
	// I know nothing about all about what pg could return as errors
	// (i thing everything except plain and clear info)
	return xerr.WrapWithType(err, errdefs.ErrTemporaryUnavailable)
}

func isPgSyntaxError(code string) bool {
	switch code {
	case errcode.SyntaxError,
		errcode.UndefinedTable,
		errcode.UndefinedParameter,
		errcode.UndefinedObject,
		errcode.DuplicateColumn,
		errcode.InvalidPassword,
		errcode.InvalidCatalogName,
		errcode.UndefinedColumn,
		errcode.UndefinedFunction,
		errcode.CardinalityViolation:
		return true
	default:
		return false
	}
}
