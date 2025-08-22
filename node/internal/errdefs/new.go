package errdefs

import (
	"github.com/go-faster/errors"
)

type option = func(e *baseError)

func New(text string, opts ...option) error {
	err := &baseError{
		err:   errors.New(text),
		stack: getTrace(2),
	}
	for _, o := range opts {
		o(err)
	}
	return err
}
