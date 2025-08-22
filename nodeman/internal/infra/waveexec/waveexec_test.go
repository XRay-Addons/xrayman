package waveexec

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWaveExec(t *testing.T) {
	fnCallsCount := 0

	fn := func(ctx context.Context) (any, error) {
		time.Sleep(1 * time.Second)
		fnCallsCount++
		return nil, nil
	}

	waveExec := NewWaveExecutor(fn)
	defer waveExec.Close()

	// run 100 calls at the same time, fnCallsCount expected to be
	// 1 (all calls in first wave) or 2 (some in the first wave, others in second)
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := waveExec.Invoke(context.Background())
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
	require.Positive(t, fnCallsCount)
	require.LessOrEqual(t, fnCallsCount, 2)
}

func TestWaveExec_Cancels(t *testing.T) {
	fnCallsCount := 0

	fn := func(ctx context.Context) (any, error) {
		time.Sleep(1 * time.Second)
		fnCallsCount++
		return nil, nil
	}

	waveExec := NewWaveExecutor(fn)
	defer waveExec.Close()

	// run 100 calls with different timeouts
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(2)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()
			_, err := waveExec.Invoke(ctx)
			assert.Error(t, err)
		}()
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
			defer cancel()
			_, err := waveExec.Invoke(ctx)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
	require.Positive(t, fnCallsCount)
	require.LessOrEqual(t, 2, fnCallsCount)
}

func TestWaveExec_EarlyClose(t *testing.T) {
	fnCallsCount := 0

	fn := func(ctx context.Context) (any, error) {
		time.Sleep(1 * time.Second)
		fnCallsCount++
		return nil, nil
	}

	waveExec := NewWaveExecutor(fn)

	// run 100 calls with early timeouts
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()
			_, err := waveExec.Invoke(ctx)
			assert.Error(t, err)
		}()
	}

	wg.Wait()
	waveExec.Close()

	require.Positive(t, fnCallsCount)
	require.LessOrEqual(t, fnCallsCount, 2)
}

func TestWaveExec_Cancel(t *testing.T) {
	fnCallsCount := 0
	fn := func(ctx context.Context) (any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.NewTimer(1 * time.Second).C:
			fnCallsCount++
			return nil, nil
		}
	}

	waveExec := NewWaveExecutor(fn)

	// run 100 calls with early timeouts
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()
			_, err := waveExec.Invoke(ctx)
			assert.Error(t, err)
		}()
	}

	wg.Wait()
	waveExec.Close()

	require.Equal(t, 0, fnCallsCount)
}
