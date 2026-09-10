package config

import (
	_ "embed"
	"net"
	"os"

	"github.com/XRay-Addons/xrayman/common/jsonval"
	"github.com/XRay-Addons/xrayman/common/xerr"
)

//go:embed server_schema.json
var serverSchema []byte

//go:embed client_schema.json
var clientSchema []byte

func Validate(c *Config) error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return xerr.InvalidArg("endpoint", c.Endpoint, err)
	}
	if err := checkDir(c.XRayDataDir); err != nil {
		return xerr.InvalidArg("xray data dir", c.XRayDataDir, err)
	}
	if err := checkFile(c.XRayServer()); err != nil {
		return xerr.InvalidArg("xray server config", c.XRayServer(), err)
	}
	if err := checkJson(c.XRayServer(), serverSchema); err != nil {
		return xerr.InvalidArg("xray server config", c.XRayServer(), err)
	}
	if err := checkFile(c.XRayClient()); err != nil {
		return xerr.InvalidArg("xray client config", c.XRayClient(), err)
	}
	if err := checkJson(c.XRayClient(), clientSchema); err != nil {
		return xerr.InvalidArg("xray client config", c.XRayClient(), err)
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

func checkJson(path string, schema []byte) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return xerr.WrapWithStack(err)
	}
	if err := jsonval.ValidateJsonData(data); err != nil {
		return xerr.Wrap(err, xerr.WithFile(path))
	}
	if err := jsonval.ValidateJsonSchema(data, schema); err != nil {
		return xerr.Wrap(err, xerr.WithFile(path))
	}
	return nil
}
