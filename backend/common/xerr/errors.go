package xerr

var (
	ErrNilArg  = Define("nil arg")
	ErrNilCall = Define("nil object call")
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
