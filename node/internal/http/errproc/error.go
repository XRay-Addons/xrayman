package errproc

import (
	"errors"
)

// error containing http response error
// and real reason for it
type Error struct {
	response *Response
	reason   error
}

func NewError(response *Response, reason error) *Error {
	return &Error{response: response, reason: reason}
}

func (e *Error) Error() string {
	return e.reason.Error()
}

func (e *Error) Is(target error) bool {
	// in C++ we call it fffuuu-pattern
	if r, ok := target.(*Response); ok {
		return e.response == r
	}
	return errors.Is(e.reason, target)
}

func (e *Error) As(target any) bool {
	if t, ok := target.(**Response); ok {
		*t = e.response
		return true
	}
	return errors.As(e.reason, target)
}

func (e *Error) Unwrap() error {
	return e.reason
}
