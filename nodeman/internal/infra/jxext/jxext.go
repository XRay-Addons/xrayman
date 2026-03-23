package jxext

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/go-faster/jx"
)

func Validate(val jx.Raw) error {
	d := jx.DecodeBytes(val)
	if err := d.Skip(); err != nil {
		return errdefs.WrapWithStack(err)
	}

	if d.Next() != jx.Invalid {
		return errdefs.New(fmt.Sprintf("validation: text after end: %v", d.Next().String()))
	}
	return nil
}
