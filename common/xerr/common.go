package xerr

var (
	ErrNilCall = New("nil object call")
	ErrNilArg  = New("nil argument passed")
)

func NilCall() error {
	return WrapWithStack(ErrNilCall)
}

func NilArg(name string) error {
	return Wrap(ErrNilArg, WithStack(), Withf("argument name: %s", name))
}
