package app

import (
	"context"
	"errors"
	"os"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestApp_New(t *testing.T) {
	logger := zaptest.NewLogger(t)
	app := New(
		WithLogger(logger),
		WithComponent("test", nil, nil),
		WithRunner("runner", func() error { return nil }, nil),
	)

	require.NotNil(t, app)
	require.Len(t, app.components, 1)
	require.Len(t, app.runners, 1)
}

func TestApp_Run_Success(t *testing.T) {
	t.Run("all operations succeed", func(t *testing.T) {
		logger := zaptest.NewLogger(t)

		// Track execution order
		execOrder := []string{}

		compInit := func() error {
			execOrder = append(execOrder, "init")
			return nil
		}
		compClose := func(ctx context.Context) error {
			execOrder = append(execOrder, "close")
			return nil
		}

		runner := func() error {
			execOrder = append(execOrder, "run")
			return nil
		}

		app := New(
			WithLogger(logger),
			WithComponent("comp1", compInit, compClose),
			WithComponent("comp2", compInit, compClose),
			WithRunner("runner1", runner, nil),
			WithRunner("runner2", runner, nil),
		)

		err := app.Run()
		require.NoError(t, err)

		// Verify execution order
		expected := []string{
			"init", "init", // Components init in order
			"run", "run", // Runners execute
			"close", "close", // Components close in reverse order
		}
		require.Equal(t, expected, execOrder)
	})
}

func TestApp_Run_InitFailure(t *testing.T) {
	logger := zaptest.NewLogger(t)

	initErr := errors.New("init error")
	initFn := func() error { return initErr }
	closeFn := func(ctx context.Context) error { return nil }
	runnerFn := func() error { return nil }

	app := New(
		WithLogger(logger),
		WithComponent("failing", initFn, closeFn),
		WithRunner("runner", runnerFn, closeFn),
	)

	err := app.Run()
	require.Error(t, err)
	require.Contains(t, err.Error(), "init app")
}

func TestApp_Run_RunnerFailure(t *testing.T) {
	logger := zaptest.NewLogger(t)

	runnerErr := errors.New("runner error")
	initFn := func() error { return nil }
	closeFn := func(ctx context.Context) error { return nil }
	runnerFn := func() error { return runnerErr }

	app := New(
		WithLogger(logger),
		WithComponent("comp", initFn, closeFn),
		WithRunner("failing", runnerFn, nil),
	)

	err := app.Run()
	require.Error(t, err)
	require.Contains(t, err.Error(), "app run")
	require.ErrorIs(t, err, runnerErr)
}

func TestApp_Close_ErrorHandling(t *testing.T) {
	logger := zaptest.NewLogger(t)

	closeErr1 := errors.New("close error 1")
	closeErr2 := errors.New("close error 2")

	app := New(
		WithLogger(logger),
		WithComponent("comp1",
			func() error { return nil },
			func(ctx context.Context) error { return closeErr1 },
		),
		WithComponent("comp2",
			func() error { return nil },
			func(ctx context.Context) error { return closeErr2 },
		),
		WithRunner("runner",
			func() error { return nil },
			nil,
		),
	)

	// Initialize components
	err := app.init()
	require.NoError(t, err)

	// Close components (should return aggregated errors)
	err = app.close()
	require.Error(t, err)
	require.Contains(t, err.Error(), "close app error")
	require.True(t, strings.Contains(err.Error(), closeErr1.Error()))
	require.True(t, strings.Contains(err.Error(), closeErr2.Error()))
}

func TestApp_Run_ContextCancel(t *testing.T) {
	logger := zaptest.NewLogger(t)
	ctx, cancel := context.WithCancel(context.Background())

	// Runner that blocks until context cancellation
	runner := func() error {
		<-ctx.Done()
		return ctx.Err()
	}

	app := New(
		WithLogger(logger),
		WithRunner("blocking", runner, nil),
	)

	// Run app in separate goroutine
	errCh := make(chan error)
	go func() {
		errCh <- app.Run()
	}()

	// Cancel context after short delay
	time.AfterFunc(10*time.Millisecond, cancel)

	// Wait for app to finish
	err := <-errCh
	require.Error(t, err)
	require.True(t, errors.Is(err, context.Canceled),
		"expected context.Canceled, got: %v", err)
}

func TestComponentLifecycle_Order(t *testing.T) {
	logger := zaptest.NewLogger(t)
	order := []string{}

	app := New(
		WithLogger(logger),
		WithComponent("first",
			func() error {
				order = append(order, "init-first")
				return nil
			},
			func(ctx context.Context) error {
				order = append(order, "close-first")
				return nil
			},
		),
		WithComponent("second",
			func() error {
				order = append(order, "init-second")
				return nil
			},
			func(ctx context.Context) error {
				order = append(order, "close-second")
				return nil
			},
		),
	)

	// Initialize
	err := app.init()
	require.NoError(t, err)

	// Close
	err = app.close()
	require.NoError(t, err)

	// Verify order
	expected := []string{
		"init-first", "init-second",
		"close-second", "close-first", // Closed in reverse order
	}
	require.Equal(t, expected, order)
}

func TestApp_Run_SignalCancel(t *testing.T) {
	t.Run("all operations succeed, close by done", func(t *testing.T) {
		logger := zaptest.NewLogger(t)

		// Track execution order
		execOrder := []string{}

		runner := func() error {
			time.Sleep(2 * time.Second)
			execOrder = append(execOrder, "run")
			return nil
		}

		app := New(
			WithLogger(logger),
			WithRunner("runner", runner, nil),
			WithSignalCancel(),
		)

		go func() {
			time.Sleep(4 * time.Second)
			p, _ := os.FindProcess(os.Getpid())
			_ = p.Signal(syscall.SIGINT)
		}()

		err := app.Run()
		require.NoError(t, err)

		// Verify execution order
		expected := []string{
			"run", // Runner execute
		}
		require.Equal(t, expected, execOrder)
	})
}

func TestApp_Run_SignalCancelInterrupt(t *testing.T) {
	t.Run("all operations succeed, close by interript", func(t *testing.T) {
		logger := zaptest.NewLogger(t)

		// Track execution order
		execOrder := []string{}

		runner := func() error {
			time.Sleep(4 * time.Second)
			execOrder = append(execOrder, "run")
			return nil
		}

		app := New(
			WithLogger(logger),
			WithRunner("runner", runner, nil),
			WithSignalCancel(),
		)

		err := app.Run()
		require.Error(t, err)

		// Verify execution order
		expected := []string{
			"run", // Runner execute
		}
		require.Equal(t, expected, execOrder)
	})
}
