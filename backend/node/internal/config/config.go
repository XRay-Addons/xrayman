package config

import (
	"path"
)

type Config struct {
	Endpoint      string
	XRayDataDir   string
	XRayConfigDir string
	PersistentDir string
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
	}

	if err := Validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
