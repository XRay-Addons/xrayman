package errdefs

func NewNilArg(name string) error {
	return Newf("nil arg: %s", name)
}

func NewNilCall() error {
	return New("nil object call")
}
