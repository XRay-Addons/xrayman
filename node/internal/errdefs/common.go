package errdefs

import "github.com/go-faster/errors"

var (
	ErrNilCall = errors.New("nil object call")
	ErrNilArg  = errors.New("nil argument passed")
	ErrUnknown = errors.New("unknown error")
)

func NewNilCall() error {
	return WrapWithStack(ErrNilCall)
}

func NewNilArg(name string) error {
	return Wrap(ErrNilArg, WithStack(), Withf("argument name: %s", name))
}

func NewUnknown(details string) error {
	return Wrap(ErrUnknown, WithStack(), Withf("details: %s", details))
}
