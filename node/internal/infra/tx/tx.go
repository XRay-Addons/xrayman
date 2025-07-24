package tx

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Tx struct {
	items           []txItem
	rollbackTimeout time.Duration
}

type Option func(*Tx)

func New(opts ...Option) Tx {
	tx := Tx{
		rollbackTimeout: defaultRollbackTimeout,
	}
	for _, o := range opts {
		o(&tx)
	}
	return tx
}

const defaultRollbackTimeout = 5 * time.Second

func WithRollbackTimeout(timeout time.Duration) Option {
	return func(tx *Tx) {
		tx.rollbackTimeout = timeout
	}
}

type Fn = func(ctx context.Context) error

func (tx *Tx) AddItem(fn, rb Fn) {
	tx.items = append(tx.items, txItem{fn: fn, rb: rb})
}

func (tx *Tx) Run(ctx context.Context) error {
	commited, err := tx.commit(ctx)
	if err == nil {
		return nil
	}
	err = fmt.Errorf("commit: %w", err)

	// Use separate context for rollback to ensure it runs
	rbCtx, cancel := context.WithTimeout(context.Background(), tx.rollbackTimeout)
	defer cancel()

	rbErrs := tx.rollback(rbCtx, commited)
	if len(rbErrs) > 0 {
		rbCombined := errors.Join(rbErrs...)
		err = errors.Join(err, fmt.Errorf("rollback: %w", rbCombined))
	}

	return fmt.Errorf("tx run: %w", err)
}

type txItem struct {
	fn Fn
	rb Fn
}

func (tx *Tx) commit(ctx context.Context) (commited int, err error) {
	for i, item := range tx.items {
		if err := item.fn(ctx); err != nil {
			return i, err
		}
	}
	return len(tx.items), nil
}

func (tx *Tx) rollback(ctx context.Context, commited int) (rbErrs []error) {
	for i := commited - 1; i >= 0; i-- {
		if tx.items[i].rb == nil {
			continue
		}
		if err := tx.items[i].rb(ctx); err != nil {
			rbErrs = append(rbErrs, err)
		}
	}
	return
}
