package config

import (
	"flag"
	"os"
	"path"
	"runtime"

	"github.com/XRay-Addons/xrayman/common/xerr"
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
	defaultDir := defaultXRayDir()
	return &Config{
		Endpoint:      "localhost:8080",
		XRayDataDir:   path.Join(defaultDir, "data"),
		XRayConfigDir: path.Join(defaultDir, "config"),
		PersistentDir: path.Join(defaultDir, "persistent"),
	}
}

func defaultXRayDir() string {
	switch runtime.GOOS {
	case "darwin", "linux":
		return "~/.local/share/xrayman/node"
	default:
		return ""
	}
}

func readCLIParams(c *Config) error {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	fs.StringVar(&c.Endpoint, "a", c.Endpoint,
		"server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22")

	fs.StringVar(&c.XRayDataDir, "d", c.XRayDataDir,
		`xray data dir, should contains geoip, geodat if routing uses it`)

	fs.StringVar(&c.XRayConfigDir, "c", c.XRayConfigDir,
		`xray configs dir, must contains xray_server.json and xray_client.json.
xray_server.json should be valid xray server config,
xray_clinet.json should be user config template,
supported template params:
  - {{ .VlessEmail }}
  - {{ .VlessUUID }}
`)

	fs.StringVar(&c.PersistentDir, "p", c.PersistentDir,
		`persistent config dir. persistent objects
(certs, secrets, config to connect to node) should be generated on-demand`)

	if err := fs.Parse(os.Args[1:]); err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func readEnvParams(c *Config) error {
	if err := env.Parse(c); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}
