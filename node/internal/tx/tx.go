package tx

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Fn = func(ctx context.Context) error

type txItem struct {
	fn Fn
	rb Fn
}

const defaultRollbackTimeout = 5 * time.Second

type Tx struct {
	items           []txItem
	RollbackTimeout time.Duration
}

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
	rollbackTimeout := defaultRollbackTimeout
	if tx.RollbackTimeout > 0 {
		rollbackTimeout = tx.RollbackTimeout
	}
	rbCtx, cancel := context.WithTimeout(context.Background(), rollbackTimeout)
	defer cancel()

	rbErrs := tx.rollback(rbCtx, commited)
	if len(rbErrs) > 0 {
		rbCombined := errors.Join(rbErrs...)
		err = errors.Join(err, fmt.Errorf("rollback: %w", rbCombined))
	}

	return fmt.Errorf("tx run: %w", err)
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
