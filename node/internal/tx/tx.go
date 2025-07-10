package tx

import (
	"errors"
	"fmt"
)

type Fn = func() error

type Tx struct {
	items []struct {
		fn Fn
		rb Fn
	}
}

func (tx *Tx) AddItem(fn, rb Fn) {
	tx.items = append(tx.items, struct {
		fn Fn
		rb Fn
	}{fn, rb})
}

func (tx *Tx) Run() error {
	commited, err := tx.commit()
	if err == nil {
		return nil
	}
	err = fmt.Errorf("commit: %w", err)

	rbErrs := tx.rollback(commited)
	if rbErrs != nil {
		err = errors.Join(err, fmt.Errorf("rollback: %w", errors.Join(rbErrs...)))
	}

	return fmt.Errorf("tx run: %w", err)
}

func (tx *Tx) commit() (commited int, err error) {
	for i, item := range tx.items {
		if err := item.fn(); err != nil {
			return i, err
		}
	}
	return len(tx.items), nil
}

func (tx *Tx) rollback(commited int) (rbErrs []error) {
	for i := commited - 1; i >= 0; i-- {
		if tx.items[i].rb == nil {
			continue
		}
		if err := tx.items[i].rb(); err != nil {
			rbErrs = append(rbErrs, err)
		}
	}
	return
}
