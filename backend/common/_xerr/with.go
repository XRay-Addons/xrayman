package xerr

import (
	"fmt"
)

func With(details string) Option {
	return func(e *xerror) {
		e.with = append(e.with, details)
	}
}

func Withf(details string, args ...any) Option {
	return func(e *xerror) {
		e.with = append(e.with, fmt.Sprintf(details, args...))
	}
}

func WithStack() Option {
	const wrappingTraceDepth = 4
	return func(e *xerror) {
		if len(e.stack) > 0 {
			return
		}
		e.stack = getTrace(wrappingTraceDepth)
	}
}

func WithType(t ErrType) Option {
	return func(e *xerror) {
		e.errtype = t
	}
}
