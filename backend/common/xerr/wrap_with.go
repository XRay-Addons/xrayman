package xerr

func WrapWith(err error, details string) error {
	return Wrap(err, With(details))
}

func WrapWithf(err error, details string, args ...any) error {
	return Wrap(err, Withf(details, args...))
}

// if stack already included, do nothing
func WrapWithStack(err error) error {
	return Wrap(err, WithStack())
}
