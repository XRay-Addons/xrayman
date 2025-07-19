package httperr

import (
	"errors"
)

// error containing http response error
// and real reason for it
type httpError struct {
	response *Response
	reason   error
}

func New(response *Response, reason error) *httpError {
	return &httpError{response: response, reason: reason}
}

func (e *httpError) Error() string {
	return e.reason.Error()
}

func (e *httpError) Is(target error) bool {
	// in C++ we call it fffuuu-pattern
	if r, ok := target.(*Response); ok {
		return e.response == r
	}
	return errors.Is(e.reason, target)
}

func (e *httpError) As(target any) bool {
	if t, ok := target.(**Response); ok {
		*t = e.response
		return true
	}
	return errors.As(e.reason, target)
}

func (e *httpError) Unwrap() error {
	return e.reason
}
