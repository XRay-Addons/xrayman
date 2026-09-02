package xerr

var (
	ErrNilArg     = Define("nil arg")
	ErrNilCall    = Define("nil object call")
	ErrInvalidArg = Define("invalid argument")
	ErrPanic      = Define("panic")
)

func NilArg(name string) error {
	return Wrap(ErrNilArg,
		WithInfof("arg name: %s", name),
		withStack(1))
}

func NilCall() error {
	return Wrap(ErrNilCall,
		withStack(1))
}

func InvalidArgf(f string, args ...any) error {
	return WrapWithInfof(ErrInvalidArg, f, args...)
}

func Panic(p any) error {
	return Wrap(ErrPanic,
		WithInfof("panic info: %v", p),
		withStack(1))
}
