package xrayapi

import (
	"errors"
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

type TxFn struct {
	Fn func() error
	Rb func() error
}

type ApiTx struct {
	fns []TxFn
}

func (tx *ApiTx) AddFn(fn TxFn) {
	tx.fns = append(tx.fns, fn)
}

func (tx *ApiTx) Execute() error {
	for idx, fn := range tx.fns {
		err := fn.Fn()
		if err == nil {
			continue
		}
		rbErr := tx.rollback(idx)
		if rbErr != nil {
			return fmt.Errorf("%w: tx function: %v, rollback: %v", errdefs.ErrXRay, err, rbErr)
		}
		return fmt.Errorf("%w: tx function: %v", errdefs.ErrXRay, err)
	}

	return nil
}

func (tx *ApiTx) rollback(failedIdx int) error {
	errs := make([]error, 0)
	for idx := failedIdx - 1; idx >= 0; idx-- {
		if err := tx.fns[idx].Rb(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
