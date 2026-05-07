package errdefs

import (
	"github.com/XRay-Addons/xrayman/nodeman/internal/infra/common/xerr"
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
