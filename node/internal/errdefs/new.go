package errdefs

import (
	"github.com/go-faster/errors"
)

type option = func(e *baseError)

func New(text string, opts ...option) error {
	const wrappingTraceDepth = 2
	err := &baseError{
		err:   errors.New(text),
		stack: getTrace(wrappingTraceDepth),
	}
	for _, o := range opts {
		o(err)
	}
	return err
}
