package uow

import "context"

type Fn[T any] = func(ctx T) error

type Storage[T any] interface {
	DoUoW(ctx context.Context, fn Fn[T]) error
}
