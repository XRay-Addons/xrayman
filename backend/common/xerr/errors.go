package xerr

var (
	ErrNilArg  = Define("nil arg")
	ErrNilCall = Define("nil object call")
	ErrPanic   = Define("panic")
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

func Panic(p any) error {
	return Wrap(ErrPanic,
		WithInfof("panic info: %v", p),
		withStack(1))
}
