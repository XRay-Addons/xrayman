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
	ErrInvaildPayload       = xerr.Define("invalid payload")
	ErrNotFound             = xerr.Define("not found")
)

func NilCall() error {
	return xerr.NilCall()
}

func NilArg(name string) error {
	return xerr.NilArg(name)
}

func AccessDenied() error {
	return xerr.Wrap(ErrAccessDenied,
		xerr.WithStack())
}

//func InvalidPayload(details string) error {
//	return xerr.Wrap(ErrInvaildPayload,
//		xerr.WithStack(),
//		xerr.WithInfo(details))
//}

/*func NotFound(details string) error {
	return xerr.Wrap(ErrNotFound,
		xerr.WithStack(),
		xerr.WithInfo(details))
}

func Connection(method, url string, code int) error {
	return xerr.Wrap(ErrConnection,
		xerr.WithStack(),
		xerr.WithInfof("%s %s: status %d", method, url, code))
}*/
