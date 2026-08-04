package config

import (
	"net"
	"os"

	"github.com/XRay-Addons/xrayman/common/jsonval"
	"github.com/XRay-Addons/xrayman/common/xerr"
)

func Validate(c *Config) error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return xerr.Wrap(err,
			xerr.WithStack(),
			xerr.WithInfof("invalid endpoint %s", c.Endpoint))
	}
	if err := checkDir(c.XRayDataDir); err != nil {
		return xerr.WrapWithInfof(err, "arg: %s", c.XRayDataDir)
	}
	if err := checkFile(c.XRayServer()); err != nil {
		return xerr.WrapWithInfo(err, "xray server cfg")
	}
	if err := checkJson(c.XRayServer()); err != nil {
		return xerr.WrapWithInfo(err, "xray server cfg")
	}
	if err := checkFile(c.XRayClient()); err != nil {
		return xerr.WrapWithInfo(err, "xray client cfg")
	}
	if err := checkJson(c.XRayClient()); err != nil {
		return xerr.WrapWithInfo(err, "xray client cfg")
	}
	// don't check c.PersistentDir, it could be created later

	return nil
}

func checkDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return xerr.WrapWithStack(err)
	}
	if !info.Mode().IsDir() {
		return xerr.Newf("file %s is not dir", path)
	}
	return nil
}

func checkFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return xerr.WrapWithStack(err)
	}
	if !info.Mode().IsRegular() {
		return xerr.Newf("file %s is not regular", path)
	}
	return nil
}

func checkJson(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return xerr.WrapWithStack(err)
	}
	if err := jsonval.ValidateJsonData(data); err != nil {
		return xerr.Wrap(err, xerr.WithFile(path))
	}
	return nil
}
