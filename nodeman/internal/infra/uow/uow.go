package uow

import "context"

type Fn[T any] func(ctx T) error

type UoW[T any] interface {
	Do(ctx context.Context, fn Fn[T]) error
}

type Storage[T any] interface {
	NewUoW() UoW[T]
}
