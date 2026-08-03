package waveexec

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testTimeS(t *testing.T, target int, dur time.Duration) {
	sec := int(dur+time.Second/2) / int(time.Second)
	require.Equal(t, target, sec)
}

func TestWaveExec(t *testing.T) {
	fnCallsCount := 0

	fn := func(ctx context.Context) (*any, error) {
		time.Sleep(1 * time.Second)
		fnCallsCount++
		return nil, nil
	}

	waveExec := New(fn)
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

	fn := func(ctx context.Context) (*any, error) {
		time.Sleep(1 * time.Second)
		fnCallsCount++
		return nil, nil
	}

	waveExec := New(fn)
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

	fn := func(ctx context.Context) (*any, error) {
		time.Sleep(1 * time.Second)
		fnCallsCount++
		return nil, nil
	}

	waveExec := New(fn)

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
	fn := func(ctx context.Context) (*any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.NewTimer(1 * time.Second).C:
			fnCallsCount++
			return nil, nil
		}
	}

	waveExec := New(fn)

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

func TestWaveExec_MultiCancel(t *testing.T) {

	fnCallsCount := 0
	fnSuccessCount := 0
	// fn lasts for 5 seconds
	fn := func(ctx context.Context) (*any, error) {
		fnCallsCount++
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.NewTimer(5 * time.Second).C:
			fnSuccessCount++
			return nil, nil
		}
	}

	waveExec := New(fn)
	var wg sync.WaitGroup

	start := time.Now()

	// 1st wave
	// call with 2 seconds timeout - should fail by timeout
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, err := waveExec.Invoke(ctx)
		testTimeS(t, 2, time.Now().Sub(start))
		assert.Error(t, err)
	}()

	time.Sleep(1 * time.Second)

	// 2nd wave - queued after 1s from start, run after 2s from start
	// should wait for sleep (1s) and its own timeout (2 sec) and fail
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, err := waveExec.Invoke(ctx)
		testTimeS(t, 3, time.Now().Sub(start))
		assert.Error(t, err)
	}()

	// should wait first call timeout (2s), second call timeout (3s) and call time (5s)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := waveExec.Invoke(ctx)
		testTimeS(t, 7, time.Now().Sub(start))
		assert.NoError(t, err)
	}()

	wg.Wait()
	waveExec.Close()

	require.Equal(t, 2, fnCallsCount)
	require.Equal(t, 1, fnSuccessCount)
	testTimeS(t, 7, time.Now().Sub(start))
}

func TestWaveExec_Timeout(t *testing.T) {
	// fn lasts for 5 seconds
	fn := func(ctx context.Context) (*any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.NewTimer(5 * time.Second).C:
			return nil, nil
		}
	}

	waveExec := New(fn)

	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := waveExec.Invoke(context.Background())
		testTimeS(t, 1, time.Now().Sub(start))
		assert.Error(t, err)
	}()

	time.Sleep(1 * time.Second)
	waveExec.Close()
	testTimeS(t, 1, time.Now().Sub(start))
	wg.Wait()
}

func TestWaveExec_Panic(t *testing.T) {
	// fn lasts for 5 seconds
	fn := func(ctx context.Context) (*any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.NewTimer(5 * time.Second).C:
			panic("funny panic")
		}
	}

	waveExec := New(fn)

	start := time.Now()

	var wg sync.WaitGroup

	// 2s timeout
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, err := waveExec.Invoke(ctx)
		testTimeS(t, 2, time.Now().Sub(start))
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	}()

	// to ensure next call in 2nd wave
	time.Sleep(10 * time.Millisecond)

	// without timeout
	_, err := waveExec.Invoke(t.Context())
	testTimeS(t, 7, time.Now().Sub(start))
	assert.ErrorIs(t, err, xerr.ErrPanic)

	waveExec.Close()
	testTimeS(t, 7, time.Now().Sub(start))
	wg.Wait()
}
