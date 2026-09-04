package xerr

var (
	ErrNilArg     = Define("nil arg")
	ErrNilCall    = Define("nil object call")
	ErrInvalidArg = Define("invalid argument")
	ErrPanic      = Define("panic")
)

func NilArg(name string) error {
	return Wrap(ErrNilArg,
		withStack(1),
		WithInfof("arg name: %s", name),
	)
}

func NilCall() error {
	return Wrap(ErrNilCall,
		withStack(1),
	)
}

func InvalidArgf(f string, args ...any) error {
	return Wrap(ErrInvalidArg,
		withStack(1),
		WithInfof(f, args...),
	)
}

func Panic(p any) error {
	return Wrap(ErrPanic,
		withStack(1),
		WithInfof("panic info: %v", p),
	)
}
