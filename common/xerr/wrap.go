package xerr

func Wrap(base error, opts ...Option) error {
	if base == nil {
		return base
	}
	err := &baseError{
		err: base,
	}
	for _, o := range opts {
		o(err)
	}
	return err
}
