package errdefs

func Wrap(base error, opts ...option) error {
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
