package fs

import (
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

func AccessFile(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errdefs.NewFileAccess(path, err)
	}
	return !info.IsDir(), nil
}

func AccessDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errdefs.NewFileAccess(path, err)
	}
	return info.IsDir(), nil
}
