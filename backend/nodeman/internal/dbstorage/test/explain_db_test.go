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

func (e *ExplainTX) Commit() error {
	err := e.tx.Commit()
	e.processExplainQuery()
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

func (e *ExplainTX) processExplainQuery() {
	defer func() { e.queries = []expquery{} }()

	if e.mode == ExplainNone {
		return
	}

	tx, err := e.db.BeginTx(context.TODO())
	if err != nil {
		e.expl.Add(fmt.Sprintf("begin explain tx error: %+v", err))
		return
	}

	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil {
			e.expl.Add(fmt.Sprintf("explain tx rollback error: %+v", err))
		}
		return
	}()

	for _, q := range e.queries {
		var b strings.Builder
		prefix := getModePrefix(e.mode)

		rows, err := tx.QueryContext(context.TODO(), prefix+q.q, q.args...)
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

		e.expl.Add(b.String())
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
