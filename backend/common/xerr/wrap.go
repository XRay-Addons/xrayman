package xerr

type option func(err error) error

func Wrap(err error, options ...option) error {
	for _, o := range options {
		err = o(err)
	}
	return err
}

func WrapWithStack(err error) error {
	err = WithStack()(err)
	return err
}

func WrapWithFile(err error, file string) error {
	return WithFile(file)(err)
}

func WrapWithInfo(err error, info string) error {
	return WithInfo(info)(err)
}

func WrapWithInfof(err error, info string, args ...any) error {
	return WithInfof(info, args...)(err)
}

func WrapWithType(err error, t error) error {
	return WithType(t)(err)
}
