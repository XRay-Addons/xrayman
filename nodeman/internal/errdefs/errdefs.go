package errdefs

import (
	"errors"
	"fmt"
	"net/url"

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

func WithFile(filename string) xerr.Option {
	return xerr.With(fmt.Sprintf("file: %s", filename))
}

func OgenErr(err error) error {
	details := ""

	var sc interface{ StatusCode() int }
	if errors.As(err, &sc) {
		details += fmt.Sprintf("Status: %d;", sc.StatusCode())
	} else {
		details += "Status: Transport error"
	}
	// add url path if exists
	var ue *url.Error
	if errors.As(err, &ue) {
		details += fmt.Sprintf("; URL: %s", ue.URL)
	}

	return xerr.WrapWith(ErrConnection, details)
}

func WrapWithFile(err error, path string) error {
	return xerr.Wrap(err, WithFile(path))
}
