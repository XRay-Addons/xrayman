package waveexec

import (
	"context"
	"sync"
)

// WaveExecutor coordinates execution of a single function (Fn) in "waves".
//
// When multiple callers invoke Run concurrently, they are grouped together,
// and only one execution of the function is performed for the entire group.
// All callers receive the same result.
//
// If a new Run call arrives after the current execution has already started,
// it will wait for the next wave and trigger a new function execution.
//
// This is useful when you want to ensure:
// - only one instance of a function runs at a time
// - concurrent calls are collapsed into a single execution
// - each wave produces fresh result (so fn invoking result is always
//   fresher than invoking time moment)
//
// Notes:
// - Context passed to Invoke uses only as canceller

type Fn = func(ctx context.Context) (any, error)

type WaveExecutor struct {
	fn       execFn
	nextWave []execWaveItem
	reqCh    chan struct{}
	mu       sync.Mutex
	done     chan struct{}
}

type execFn = func(context.Context) execResult

type execResult struct {
	result any
	err    error
}

type execWaveItem struct {
	ctx    context.Context
	result chan execResult
}

func NewWaveExecutor(fn Fn) *WaveExecutor {
	we := &WaveExecutor{
		fn: func(ctx context.Context) execResult {
			res, err := fn(ctx)
			return execResult{result: res, err: err}
		},
		reqCh: make(chan struct{}, 1),
		done:  make(chan struct{}),
	}
	go we.runExecLoop()
	return we
}

func (we *WaveExecutor) Close() {
	close(we.reqCh)
	<-we.done
}

func (we *WaveExecutor) Invoke(ctx context.Context) (any, error) {
	waveItem := execWaveItem{
		ctx:    ctx,
		result: make(chan execResult, 1),
	}

	we.mu.Lock()
	we.nextWave = append(we.nextWave, waveItem)
	we.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case we.reqCh <- struct{}{}:
		// schedule next wave, it contains current call
	default:
		// if reqCh is full, next wave is already
		// scheduled, current call is already in the queue
	}

	select {
	case res := <-waveItem.result:
		return res.result, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (we *WaveExecutor) runExecLoop() {
	defer close(we.done)
	for range we.reqCh {
		we.mu.Lock()
		currWave := we.nextWave
		we.nextWave = nil
		we.mu.Unlock()

		if len(currWave) == 0 {
			// nothing to execute
			continue
		}

		// create merged ctx from all calls ctxs, which
		// lasts till at least one ctx is alive. it's ok because
		// Invoke waiting for result till its own passed ctx alive
		ctx := we.anyAliveContext(currWave...)

		res := we.fn(ctx)
		for _, item := range currWave {
			select {
			case item.result <- res:
			default:
			}
			close(item.result)
		}
	}
}

func (we *WaveExecutor) anyAliveContext(items ...execWaveItem) context.Context {
	if len(items) == 0 {
		return context.Background()
	}
	if len(items) == 1 {
		return items[0].ctx
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	for _, item := range items {
		wg.Add(1)
		go func() {
			select {
			case <-item.ctx.Done():
				wg.Done()
			case <-ctx.Done():
				// merged ctx cancelled outside
			}
		}()
	}

	go func() {
		wg.Wait()
		cancel()
	}()

	return ctx
}
