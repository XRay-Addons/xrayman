package errdefs

import (
	"github.com/XRay-Addons/xrayman/common/xerr"
)

var (
	ErrNilCall              = xerr.ErrNilArg
	ErrNilArg               = xerr.ErrNilArg
	ErrConnection           = xerr.Define("connection")
	ErrTemporaryUnavailable = xerr.Define("temporary unavailable")
	ErrAccessDenied         = xerr.Define("access denied")
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
