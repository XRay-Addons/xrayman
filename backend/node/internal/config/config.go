package config

import "path"

type Config struct {
	Endpoint      string `env:"ENDPOINT"`
	XRayDataDir   string `env:"XRAY_DATA_DIR"`
	XRayConfigDir string `env:"XRAY_CONFIG_DIR"`
	PersistentDir string `env:"PERSISTENT_DIR"`
}

func (c *Config) XRayServer() string {
	return path.Join(c.XRayConfigDir, "xray_server.json")
}

func (c *Config) XRayClient() string {
	return path.Join(c.XRayConfigDir, "xray_client.json")
}
