package waveexec

import (
	"sync"
)

// queue to push items one by one and pull
// all collected items via channel
type queue[T any] struct {
	Ready chan struct{}
	items []T
	mu    sync.Mutex
}

func makeQueue[T any]() queue[T] {
	return queue[T]{
		Ready: make(chan struct{}, 1),
	}
}

func (q *queue[T]) Push(item T) {
	q.mu.Lock()
	q.items = append(q.items, item)
	if len(q.items) == 1 {
		q.Ready <- struct{}{}
	}
	q.mu.Unlock()
}

func (q *queue[T]) Next() []T {
	q.mu.Lock()
	defer func() {
		q.items = nil
		q.mu.Unlock()
	}()
	return q.items
}
