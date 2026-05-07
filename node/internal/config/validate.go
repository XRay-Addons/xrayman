package config

import (
	"net"
	"os"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

func Validate(c Config) error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return xerr.Wrap(err, xerr.WithStack(),
			xerr.Withf("invalid endpoint %s", c.Endpoint))
	}
	if err := checkDir(c.XRayData()); err != nil {
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

func checkFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return xerr.WrapWithStack(err)
	}
	if !info.Mode().IsRegular() {
		return xerr.New("file is not regular", errdefs.WithFile(path))
	}
	return nil
}

func checkDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return xerr.WrapWithStack(err)
	}
	if !info.Mode().IsDir() {
		return xerr.New("file is not dir", errdefs.WithFile(path))
	}
	return nil
}
