package config

import (
	"net"
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

func Validate(c Config) error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return errdefs.Wrap(err, errdefs.WithStack(),
			errdefs.Withf("invalid endpoint %s", c.Endpoint))
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
		return errdefs.Wrap(err, errdefs.WithStack(), errdefs.WithFile(path))
	}
	if !info.Mode().IsRegular() {
		return errdefs.New("exectuable is not regular file",
			errdefs.WithFile(path))
	}
	perm := info.Mode().Perm()
	if perm&0111 != 0 {
		return nil
	}
	return errdefs.New("file is not executable for current user",
		errdefs.WithFile(path))
}

func checkFile(path string) error {
	exists, err := checkFileExists(path)
	if err != nil {
		return errdefs.Wrap(err, errdefs.WithStack(), errdefs.WithFile(path))
	}
	if !exists {
		return errdefs.New("file not exists", errdefs.WithFile(path))
	}
	return nil
}

func checkFileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	if !info.Mode().IsRegular() {
		return false, errdefs.New("file is not regular", errdefs.WithFile(path))
	}
	return true, nil
}
