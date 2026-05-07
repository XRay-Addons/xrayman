package fs

import (
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/common/xerr"
)

func AccessFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, xerr.Wrap(err, xerr.WithStack(), errdefs.WithFile(path))
	}
	return !info.IsDir(), nil
}

func AccessDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, xerr.Wrap(err, xerr.WithStack(), errdefs.WithFile(path))
	}
	return info.IsDir(), nil
}
