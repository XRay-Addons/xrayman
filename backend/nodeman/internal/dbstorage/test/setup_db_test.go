package dbstoragetest

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

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

func setupTestDB(t *testing.T, logger *zap.Logger) (
	storage *dbstorage.Storage,
	db *ExplainDB,
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
		"postgres://test:test@%s:%s/testdb?sslmode=disable",
		host,
		port.Port(),
	)

	sqldb, err := sql.Open("pgx", connStr)
	require.NoError(t, err)

	pingCtx, cancelPing := context.WithTimeout(ctx, 5*time.Second)
	defer cancelPing()

	err = sqldb.PingContext(pingCtx)
	require.NoError(t, err, "DB Running, ping failed")

	db = &ExplainDB{
		db:   sqldb,
		mode: ExplainNone,
	}
	storage, err = dbstorage.New(db, dbstorage.WithLogger(logger))
	require.NoError(t, err)
	storage.Migrate(ctx)

	t.Cleanup(func() {
		_ = db.Close()
	})

	return
}
