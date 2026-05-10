package errdefs

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/common/xerr"
)

var (
	ErrNilCall              = xerr.ErrNilArg
	ErrNilArg               = xerr.ErrNilArg
	ErrConnection           = xerr.New("connection")
	ErrTemporaryUnavailable = xerr.New("temporary unavailable")
	ErrAccessDenied         = xerr.New("access denied")
)

func NilCall() error {
	return xerr.NilCall()
}

func NilArg(name string) error {
	return xerr.NilArg(name)
}

func AccessDenied() error {
	return xerr.WrapWithStack(ErrAccessDenied)
}

func WithFile(filename string) xerr.Option {
	return xerr.With(fmt.Sprintf("file: %s", filename))
}

func WrapWithFile(err error, path string) error {
	return xerr.Wrap(err, WithFile(path))
}
