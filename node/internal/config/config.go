package config

import "path"

type Config struct {
	Endpoint  string `env:"ENDPOINT"`
	AccessKey string `env:"ACCESS_KEY"`

	XRayDir string `env:"XRAY_DIR"`
}

func (c *Config) XRayExec() string {
	return path.Join(c.XRayDir, "xray")
}

func (c *Config) XRayServer() string {
	return path.Join(c.XRayDir, "xray_server.json")
}

func (c *Config) XRayClient() string {
	return path.Join(c.XRayDir, "xray_client.json")
}
