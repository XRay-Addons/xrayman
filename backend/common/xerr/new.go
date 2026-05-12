package xerr

import (
	"fmt"

	"github.com/go-faster/errors"
)

type Option = func(e *xerror)

func New(text string, opts ...Option) error {
	const wrappingTraceDepth = 2
	err := &xerror{
		err:   errors.New(text),
		stack: getTrace(wrappingTraceDepth),
	}
	for _, o := range opts {
		o(err)
	}
	return err
}

func Newf(format string, a ...any) error {
	const wrappingTraceDepth = 2
	err := &xerror{
		err:   errors.New(fmt.Sprintf(format, a...)),
		stack: getTrace(wrappingTraceDepth),
	}

	return err
}
