package config

import "path"

type Config struct {
	Endpoint      string `env:"ENDPOINT"`
	XRayDir       string `env:"XRAY_DIR"`
	PersistentDir string `env:"PERSISTENT_DIR"`
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

/*func (c *Config) HasCerts() bool {
	return c.nodeCrt != "" ||
		c.nodeKey != "" ||
		c.rootCrt != ""
}

func (c *Config) NodeCrt() string {
	return c.nodeCrt
}

func (c *Config) NodeKey() string {
	return c.nodeKey
}

func (c *Config) RootCrt() string {
	return c.rootCrt
}*/
