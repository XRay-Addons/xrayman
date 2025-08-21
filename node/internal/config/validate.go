package config

import (
	"net"
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

func Validate(c Config) error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return errdefs.Withf(errdefs.WithStack(err),
			"invalid endpoint %s", c.Endpoint)
	}
	if err := checkExecutable(c.XRayExec()); err != nil {
		return err
	}
	if err := checkFile(c.XRayServer()); err != nil {
		return err
	}
	if err := checkFile(c.XRayClient()); err != nil {
		return err
	}

	return nil
}

func checkExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return errdefs.Withf(errdefs.WithStack(err),
			"invalid executable %s", path)
	}
	if !info.Mode().IsRegular() {
		return errdefs.Newf("exectuable %s is not regular file", path)
	}
	perm := info.Mode().Perm()
	if perm&0111 != 0 {
		return nil
	}
	return errdefs.Newf("file %ы executable for current user", path)
}

func checkFile(path string) error {
	exists, err := checkFileExists(path)
	if err != nil {
		return errdefs.WithStack(err)
	}
	if !exists {
		return errdefs.Newf("file not exists: %v", path)
	}
	return nil
}

func checkFileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	if !info.Mode().IsRegular() {
		return false, errdefs.Newf("file %s is not regular file", path)
	}
	return true, nil
}
