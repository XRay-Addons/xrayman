package supervisor

import (
	"context"
	"sync"
	"time"
)

type Supervisor struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func New() *Supervisor {
	ctx, cancel := context.WithCancel(context.Background())

	return &Supervisor{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Supervisor) Close() {
	if s == nil {
		return
	}
	s.cancel()
	s.wg.Wait()
	return
}

func (s *Supervisor) Go(fn func(context.Context), timeout time.Duration) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()
		ctx, cancel := context.WithTimeout(s.ctx, timeout)
		defer cancel()
		fn(ctx)
	}()
}
