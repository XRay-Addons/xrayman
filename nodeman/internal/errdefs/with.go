package errdefs

import "fmt"

func With(details string) option {
	return func(e *baseError) {
		e.with = append(e.with, details)
	}
}

func Withf(details string, args ...any) option {
	return func(e *baseError) {
		e.with = append(e.with, fmt.Sprintf(details, args...))
	}
}

func WithStack() option {
	const wrappingTraceDepth = 3
	return func(e *baseError) {
		e.stack = getTrace(wrappingTraceDepth)
	}
}

func WithoutStack() option {
	return func(e *baseError) {
		e.stack = []string{}
	}
}

func WithFile(path string) option {
	return Withf("filepath: %s", path)
}
