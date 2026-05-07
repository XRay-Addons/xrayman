package errdefs

import (
	"errors"
	"fmt"
	"net/url"
)

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
	const wrappingTraceDepth = 4
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

func WithOgen() option {
	return func(e *baseError) {
		// add status code if exists
		e.with = append(e.with, e.err.Error())
		var sc interface{ StatusCode() int }
		if errors.As(e.err, &sc) {
			e.with = append(e.with, fmt.Sprintf("Status: %d", sc.StatusCode()))
		} else {
			e.with = append(e.with, "Status: Transport error")
		}
		// add url path if exists
		var ue *url.Error
		if errors.As(e.err, &ue) {
			e.with = append(e.with, fmt.Sprintf("URL: %s", ue.URL))
		}
		// replace error
		e.err = ErrConnection
	}
}
