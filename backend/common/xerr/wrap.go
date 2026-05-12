package xerr

import "errors"

func Wrap(base error, opts ...Option) error {
	if base == nil {
		return base
	}

	var err *xerror
	if !errors.As(base, &err) {
		err = &xerror{err: base}
	}

	for _, o := range opts {
		o(err)
	}
	return err
}
