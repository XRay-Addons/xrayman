package errdefs

import (
	"github.com/XRay-Addons/xrayman/common/xerr"
)

var (
	ErrNilCall              = xerr.ErrNilArg
	ErrNilArg               = xerr.ErrNilArg
	ErrConnection           = xerr.New("connection")
	ErrTemporaryUnavailable = xerr.New("temporary unavailable")
	ErrAccessDenied         = xerr.New("access denied")
	ErrInvaildPayload       = xerr.New("invalid payload")
	ErrNotFound             = xerr.New("not found")
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

func InvalidPayload(details string) error {
	return xerr.WrapWith(ErrInvaildPayload, details)
}

func NotFound(details string) error {
	return xerr.WrapWith(ErrNotFound, details)
}
