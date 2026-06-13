package xerr

func Wrap(base error, opts ...Option) error {
	if base == nil {
		return base
	}

	var err *xerror
	if xbase, ok := base.(*xerror); ok {
		err = copyXErr(xbase)
	} else {
		err = &xerror{err: base}
	}

	for _, o := range opts {
		o(err)
	}

	return err
}

func copyXErr(e *xerror) *xerror {
	return &xerror{
		err:     e.err,
		errtype: e.errtype,
		with:    append([]string{}, e.with...),
		stack:   append([]string{}, e.stack...),
	}
}
