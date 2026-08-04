package config

import (
	"path"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Endpoint      string
	XRayDataDir   string
	XRayConfigDir string
	PersistentDir string
	LogLevel      zapcore.Level
}

func (c *Config) XRayServer() string {
	return path.Join(c.XRayConfigDir, "xray_server.json")
}

func (c *Config) XRayClient() string {
	return path.Join(c.XRayConfigDir, "xray_client.json")
}

func NewConfig(cli *CLI) (*Config, error) {
	cfg := &Config{
		Endpoint:      cli.Endpoint,
		XRayDataDir:   cli.XRayDataDir,
		XRayConfigDir: cli.XRayConfigDir,
		PersistentDir: cli.PersistentDir,
		LogLevel:      cli.LogLevel,
	}

	if err := Validate(cfg); err != nil {
		return nil, xerr.WrapWithInfo(err, "validate config")
	}

	return cfg, nil
}
