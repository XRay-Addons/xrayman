package errdefs

import (
	"fmt"

	"github.com/go-faster/errors"
)

func With(err error, w string) error {
	return &errWith{
		err:  err,
		with: w,
	}
}

func Withf(err error, w string, args ...any) error {
	return &errWith{
		err:  err,
		with: fmt.Sprintf(w, args...),
	}
}

type errWith struct {
	err  error
	with string
}

var _ error = (*errWith)(nil)
var _ errors.Wrapper = (*errWith)(nil)

func (w *errWith) Error() string {
	if w.with == "" {
		return w.err.Error()
	}
	return fmt.Sprintf("%s\nwith %s",
		w.err.Error(), w.with)
}

func (e *errWith) Unwrap() error {
	return e.err
}
