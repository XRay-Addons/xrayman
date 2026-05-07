package xerr

import (
	"fmt"
)

func With(details string) Option {
	return func(e *baseError) {
		e.with = append(e.with, details)
	}
}

func Withf(details string, args ...any) Option {
	return func(e *baseError) {
		e.with = append(e.with, fmt.Sprintf(details, args...))
	}
}

func WithStack() Option {
	const wrappingTraceDepth = 4
	return func(e *baseError) {
		e.stack = getTrace(wrappingTraceDepth)
	}
}

func WithoutStack() Option {
	return func(e *baseError) {
		e.stack = []string{}
	}
}
