package config

import (
	"flag"
	"os"

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
	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		XRayExecPath:         "/usr/local/bin/xrayman/xray",
		XRayServerConfigPath: "/var/etc/xray/server.yaml",
		XRayClientConfigPath: "/var/etc/xray/client.yaml",
		Endpoint:             "localhost:8080",
	}
}

func readCLIParams(c *Config) error {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	fs.StringVar(&c.XRayExecPath, "x", c.XRayExecPath,
		"xray executable path")

	fs.StringVar(&c.XRayServerConfigPath, "s", c.XRayServerConfigPath,
		"xray server config path")

	fs.StringVar(&c.XRayClientConfigPath, "c", c.XRayClientConfigPath,
		"xray client config path")

	fs.StringVar(&c.Endpoint, "a", c.Endpoint,
		"server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22")

	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	return nil
}

func readEnvParams(c *Config) error {
	return env.Parse(c)
}
