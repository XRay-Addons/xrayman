package jxext

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/go-faster/jx"
)

func Validate(val jx.Raw) error {
	d := jx.DecodeBytes(val)
	if err := d.Skip(); err != nil {
		return xerr.WrapWithStack(err)
	}

	if d.Next() != jx.Invalid {
		return xerr.New(fmt.Sprintf("validation: text after end: %v", d.Next().String()))
	}
	return nil
}
