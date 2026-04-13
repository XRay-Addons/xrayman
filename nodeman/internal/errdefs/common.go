package errdefs

import "github.com/go-faster/errors"

var (
	ErrNilCall = errors.New("nil object call")
	ErrNilArg  = errors.New("nil argument passed")
	ErrConnection = errors.New("connection")
)

func NewNilCall() error {
	return WrapWithStack(ErrNilCall)
}

func NewNilArg(name string) error {
	return Wrap(ErrNilArg, WithStack(), Withf("argument name: %s", name))
}

