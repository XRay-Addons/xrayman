package dbstoragetest

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage"
)

type ExplainDB struct {
	db   *sql.DB
	expl Explanations
	mode ExplainMode
}

type ExplainTX struct {
	tx      dbstorage.TX
	db      dbstorage.DB
	queries []expquery
	mode    ExplainMode
	expl    *Explanations
}

type expquery struct {
	q    string
	args []any
}

var _ dbstorage.DB = (*ExplainDB)(nil)
var _ dbstorage.TX = (*ExplainTX)(nil)

// //////////////////////////////////////////////////////////////////////////////
// ExplainDB impl
func (e *ExplainDB) Raw() *sql.DB {
	return e.db
}

func (e *ExplainDB) WithExplanations(name string, m ExplainMode) *Explanations {
	e.mode = m
	e.expl = Explanations{
		mode:  m,
		name:  name,
		start: time.Now(),
	}
	return &e.expl
}

func (e *ExplainDB) BeginTx(ctx context.Context) (dbstorage.TX, error) {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &ExplainTX{
		tx:   tx,
		db:   e,
		mode: e.mode,
		expl: &e.expl,
	}, nil
}

func (e *ExplainDB) Close() error {
	if e == nil || e.db == nil {
		return nil
	}
	if err := e.db.Close(); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}

func (e *ExplainDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	e.processExplainQuery(ctx, query, args...)
	return e.db.ExecContext(ctx, query, args...)
}

func (e *ExplainDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	e.processExplainQuery(ctx, query)
	return e.db.PrepareContext(ctx, query)
}

func (e *ExplainDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	e.processExplainQuery(ctx, query, args...)
	return e.db.QueryContext(ctx, query, args...)
}

func (e *ExplainDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	e.processExplainQuery(ctx, query, args...)
	return e.db.QueryRowContext(ctx, query, args...)
}

func (e *ExplainDB) processExplainQuery(ctx context.Context, query string, args ...any) {
	processExplainQuery(ctx, e, []expquery{
		{q: query, args: args},
	}, e.mode, &e.expl)
}

// //////////////////////////////////////////////////////////////////////////////
// ExplainTX impl
func (e *ExplainTX) Commit() error {
	err := e.tx.Commit()
	e.processExplainQuery(context.TODO())
	return err
}

func (e *ExplainTX) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return e.tx.PrepareContext(ctx, query)
}

// Rollback implements TX.
func (e *ExplainTX) Rollback() error {
	e.queries = []expquery{}
	return e.tx.Rollback()
}

func (e *ExplainTX) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	e.queries = append(e.queries, expquery{q: query, args: args})
	return e.tx.ExecContext(ctx, query, args...)
}

func (e *ExplainTX) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	e.queries = append(e.queries, expquery{q: query, args: args})
	return e.tx.QueryContext(ctx, query, args...)

}

func (e *ExplainTX) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	e.queries = append(e.queries, expquery{q: query, args: args})
	return e.tx.QueryRowContext(ctx, query, args...)
}

func (e *ExplainTX) processExplainQuery(ctx context.Context) {
	defer func() { e.queries = []expquery{} }()
	processExplainQuery(ctx, e.db, e.queries, e.mode, e.expl)
}

// //////////////////////////////////////////////////////////////////////////////
// explanation impl
func processExplainQuery(ctx context.Context, db dbstorage.DB, queries []expquery, mode ExplainMode, expl *Explanations) {
	if mode == ExplainNone {
		return
	}

	tx, err := db.BeginTx(context.TODO())
	if err != nil {
		expl.Add(fmt.Sprintf("begin explain tx error: %+v", err))
		return
	}

	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil {
			expl.Add(fmt.Sprintf("explain tx rollback error: %+v", err))
		}
		return
	}()

	for _, q := range queries {
		var b strings.Builder
		prefix := getModePrefix(mode)

		rows, err := tx.QueryContext(ctx, prefix+q.q, q.args...)
		if err != nil {
			fmt.Fprintf(&b, "explain error: %+v", err)
			return
		}
		defer func() {
			if err := rows.Close(); err != nil {
				fmt.Fprintf(&b, "explain close rows error: %+v", err)
			}
		}()

		for rows.Next() {
			var line string
			rows.Scan(&line)
			fmt.Fprintln(&b, line)
		}

		expl.Add(b.String())
	}
}

func getModePrefix(m ExplainMode) string {
	switch m {
	case ExplainText:
		return "EXPLAIN "
	case ExplainJson:
		return "EXPLAIN (FORMAT JSON) "
	case ExplainAnalyze:
		return "EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) "
	default:
		return ""
	}
}
