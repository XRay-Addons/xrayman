package config

import (
	"flag"
	"os"
	"path"
	"runtime"

	"github.com/caarlos0/env/v6"
)

func LoadConfig() (*Config, error) {
	cfg := defaultConfig()
	if err := readCLIParams(cfg); err != nil {
		return nil, err
	}
	if err := readEnvParams(cfg); err != nil {
		return nil, err
	}
	initCertPaths(cfg)

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		Endpoint:  "localhost:8080",
		AccessKey: "",
		XRayDir:   defaultXRayManDir(),
	}
}

func defaultXRayManDir() string {
	switch runtime.GOOS {
	case "darwin", "linux":
		return "/usr/local/bin/xrayman"
	default:
		return ""
	}
}

func readCLIParams(c *Config) error {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	fs.StringVar(&c.Endpoint, "a", c.Endpoint,
		"server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22")
	fs.StringVar(&c.AccessKey, "k", c.AccessKey,
		"key to access to this node")
	fs.StringVar(&c.XRayDir, "x", c.XRayDir,
		`xray binaries and configs dir. must contains
xray, xray_server.json, xray_client.json. to encrypt connection
add certificates node.crt node.key ca.crt`)

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	return nil
}

func readEnvParams(c *Config) error {
	return env.Parse(c)
}

func initCertPaths(c *Config) {
	{
		nodeCrt := path.Join(c.XRayDir, "node.crt")
		if exists, err := checkFileExists(nodeCrt); err == nil && exists {
			c.nodeCrt = nodeCrt
		}
	}
	{
		nodeKey := path.Join(c.XRayDir, "node.key")
		if exists, err := checkFileExists(nodeKey); err == nil && exists {
			c.nodeKey = nodeKey
		}
	}
	{
		rootCrt := path.Join(c.XRayDir, "ca.crt")
		if exists, err := checkFileExists(rootCrt); err == nil && exists {
			c.rootCrt = rootCrt
		}
	}
}
