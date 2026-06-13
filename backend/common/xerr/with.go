package xerr

import (
	"errors"
	"fmt"
)

func WithDetails[T any](details T) option {
	return func(err error) error {
		return &withDetails[T]{
			details: details,
			err:     err,
		}
	}
}

func ExtractDetails[T any](err error) *T {
	var w *withDetails[T]
	if errors.As(err, &w) {
		return &(w.details)
	}
	return nil
}

func WithStack() option {
	return withStack(0)
}

func WithFile(filename string) option {
	return WithInfof("file: %s", filename)
}

func WithInfo(info string) option {
	return WithDetails(info)
}

func WithInfof(infof string, args ...any) option {
	return WithDetails(fmt.Sprintf(infof, args...))
}

func WithType(t error) option {
	return func(err error) error {
		return &withType{
			withDetails: withDetails[error]{
				err:     err,
				details: t,
			},
		}
	}
}

type withDetails[T any] struct {
	details T
	err     error
}

func (w *withDetails[T]) Error() string {
	return w.err.Error()
}

func (w *withDetails[T]) Details() string {
	return fmt.Sprintf("%T: %v", w.details, w.details)
}

func (w *withDetails[T]) Unwrap() error {
	if w == nil {
		return nil
	}
	return w.err
}

func (w *withDetails[T]) Format(f fmt.State, verb rune) {
	if w == nil || w.err == nil {
		return
	}

	// add details
	if f.Flag('+') {
		fmt.Fprint(f, Render(w))
	} else {
		fmt.Fprint(f, w.Error())
	}
}

func withStack(skip int) option {
	return func(err error) error {
		// don't add more than one trace
		if trace := ExtractDetails[*trace](err); trace != nil {
			return err
		}
		trace := getTrace(3 + skip)
		err = WithDetails(&trace)(err)
		return err
	}
}

type withType struct {
	withDetails[error]
}

func (w *withType) Details() string {
	return fmt.Sprintf("type: %v", w.details)
}

func (w *withType) Is(target error) bool {
	return errors.Is(w.details, target) ||
		errors.Is(w.err, target)
}

func (w *withType) As(target any) bool {
	return errors.As(w.details, target) ||
		errors.As(w.err, target)
}
