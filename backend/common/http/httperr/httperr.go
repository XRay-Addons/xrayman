package httperr

import "errors"

// wrap error with http error. all handlers errors passed throutghtought
// NewError(). This method log real detailed error and extract and return
// api.ErrorStatusCode

func WithStatus[ApiStatus any](err error, s *ApiStatus) error {
	return &wrapping[ApiStatus]{err: err, s: s}
}

func ExtractStatus[ApiStatus any](err error) (error, *ApiStatus) {
	var s status[ApiStatus]
	if errors.As(err, &s) {
		return s.Unwrap(), s.Status()
	}

	return err, nil
}

type status[ApiStatus any] interface {
	Unwrap() error
	Status() *ApiStatus
}

type wrapping[ApiStatus any] struct {
	err error
	s   *ApiStatus
}

var _ error = (*wrapping[any])(nil)

func (w *wrapping[ApiStatus]) Error() string {
	return w.err.Error()
}

func (w *wrapping[ApiStatus]) Unwrap() error {
	return w.err
}

func (w *wrapping[ApiStatus]) Status() *ApiStatus {
	return w.s
}
