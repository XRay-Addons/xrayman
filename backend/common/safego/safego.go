package safego

import "github.com/XRay-Addons/xrayman/common/xerr"

func Invoke(fn func() error) (err error) {
	defer func() {
		// panic to error
		if p := recover(); p != nil {
			err = xerr.Panic(p)
		}
	}()
	return fn()
}
