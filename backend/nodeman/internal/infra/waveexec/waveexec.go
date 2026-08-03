package waveexec

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
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

type Fn[T any] = func(ctx context.Context) (*T, error)

type WaveExecutor[T any] struct {
	fn Fn[T]

	// waves of execution requests
	q queue[execWaveItem[T]]

	// channel to cancel request
	cancelled chan struct{}

	// channel for notify all done
	done chan struct{}
}

type execFn[T any] = func(context.Context) execResult[T]

type execResult[T any] struct {
	result *T
	err    error
}

type execWaveItem[T any] struct {
	ctx    context.Context
	result chan execResult[T]
}

func New[T any](fn Fn[T]) *WaveExecutor[T] {
	we := &WaveExecutor[T]{
		fn:        fn,
		q:         makeQueue[execWaveItem[T]](),
		cancelled: make(chan struct{}),
		done:      make(chan struct{}),
	}
	go func() {
		we.runExecLoop()
		close(we.done)
	}()
	return we
}

func (we *WaveExecutor[T]) Close() {
	close(we.cancelled)
	<-we.done
}

func (we *WaveExecutor[T]) Invoke(ctx context.Context) (*T, error) {
	item := execWaveItem[T]{
		ctx:    ctx,
		result: make(chan execResult[T], 1),
	}
	we.q.Push(item)

	select {
	case res := <-item.result:
		return res.result, res.err
	case <-ctx.Done():
		return nil, xerr.WrapWithStack(ctx.Err())
	case <-we.cancelled:
		return nil, xerr.WrapWithStack(context.Canceled)
	}
}

func (we *WaveExecutor[T]) runExecLoop() {
	for {
		select {
		case <-we.cancelled:
			return
		case <-we.q.Ready:
			items := we.q.Next()
			ctx := we.anyAliveContext(items)
			res, err := we.safeFn(ctx)
			for i := range items {
				items[i].result <- execResult[T]{result: res, err: err}
			}
		}
	}
}

func (we *WaveExecutor[T]) anyAliveContext(items []execWaveItem[T]) context.Context {
	if len(items) == 0 {
		return context.Background()
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()

		for _, item := range items {
			select {
			case <-item.ctx.Done():
				continue
			case <-we.cancelled:
				return
			}
		}
		return
	}()

	return ctx
}

func (we *WaveExecutor[T]) safeFn(ctx context.Context) (res *T, err error) {
	defer func() {
		if p := recover(); p != nil {
			res = nil
			err = xerr.Panic(p)
		}
	}()
	res, err = we.fn(ctx)
	return
}
