package jsonval

import (
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/go-faster/jx"
)

func ValidateJsonData(data []byte) error {
	d := jx.DecodeBytes(data)
	if err := d.Skip(); err != nil {
		return xerr.WrapWithStack(err)
	}

	if d.Next() != jx.Invalid {
		return xerr.Newf("validation: text after end: %v", d.Next().String())
	}
	return nil
}
