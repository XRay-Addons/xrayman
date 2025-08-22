package errdefs

func WrapWith(err error, details string) error {
	return Wrap(err, With(details))
}

func WrapWithf(err error, details string, args ...any) error {
	return Wrap(err, Withf(details, args...))
}

func WrapWithStack(err error) error {
	return Wrap(err, WithStack())
}

func WrapWithFile(err error, path string) error {
	return Wrap(err, WithFile(path))
}
