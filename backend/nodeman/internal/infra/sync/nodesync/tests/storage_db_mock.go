package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/dbstorage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

// //////////////////////////////////////////////////////////////////////////////
// stub logger for testcontainers
type noopLogConsumer struct{}

func (c *noopLogConsumer) Accept(l testcontainers.Log) {
	// Intentionally discard logs
}

func (c *noopLogConsumer) Printf(format string, v ...any) {
	// do nothing
}

type TestDB struct {
	db *sql.DB
}

type TestTX struct {
	tx dbstorage.TX
	db dbstorage.DB
}

var _ dbstorage.DB = (*TestDB)(nil)
var _ dbstorage.TX = (*TestTX)(nil)

// //////////////////////////////////////////////////////////////////////////////
// TestDB impl
func (e *TestDB) Raw() *sql.DB {
	return e.db
}

func (e *TestDB) BeginTx(ctx context.Context) (dbstorage.TX, error) {
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &TestTX{
		tx: tx,
		db: e,
	}, nil
}

func (e *TestDB) Close() error {
	if e == nil || e.db == nil {
		return nil
	}
	if err := e.db.Close(); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}

func (e *TestDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return e.db.ExecContext(ctx, query, args...)
}

func (e *TestDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return e.db.PrepareContext(ctx, query)
}

func (e *TestDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return e.db.QueryContext(ctx, query, args...)
}

func (e *TestDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return e.db.QueryRowContext(ctx, query, args...)
}

// //////////////////////////////////////////////////////////////////////////////
// TestTX impl
func (e *TestTX) Commit() error {
	err := e.tx.Commit()
	return err
}

func (e *TestTX) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return e.tx.PrepareContext(ctx, query)
}

// Rollback implements TX.
func (e *TestTX) Rollback() error {
	return e.tx.Rollback()
}

func (e *TestTX) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return e.tx.ExecContext(ctx, query, args...)
}

func (e *TestTX) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return e.tx.QueryContext(ctx, query, args...)

}

func (e *TestTX) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return e.tx.QueryRowContext(ctx, query, args...)
}

// db setup
func setupTestDB(t *testing.T, logger *zap.Logger) (
	storage *dbstorage.Storage,
	db *TestDB,
) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)

	req := testcontainers.ContainerRequest{
		Image: "postgres:17",
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
		},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30 * time.Second),
		LogConsumerCfg: &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{
				&noopLogConsumer{},
			},
		},
	}

	container, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
			Logger:           &noopLogConsumer{},
		},
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = container.Terminate(context.Background())
	})

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connStr := fmt.Sprintf(
		"postgresql://test:test@%s:%s/testdb?sslmode=disable",
		host,
		port.Port(),
	)

	sqldb, err := sql.Open("pgx", connStr)
	require.NoError(t, err)

	pingCtx, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()

	err = sqldb.PingContext(pingCtx)
	require.NoError(t, err, "DB Running, ping failed")

	db = &TestDB{
		db: sqldb,
	}
	storage, err = dbstorage.New(db, dbstorage.WithLogger(logger))
	require.NoError(t, err)
	storage.Migrate(ctx)

	t.Cleanup(func() {
		_ = db.Close()
	})

	return
}
