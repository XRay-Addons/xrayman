package errdefs

import xerr "github.com/XRay-Addons/xrayman/common/xerr"

func PayloadErr(err error) error {
	return xerr.Wrap(err, xerr.WithStack(), xerr.WithType(ErrInvaildPayload))
}
