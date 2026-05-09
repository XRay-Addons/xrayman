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
	initCertPaths(cfg)

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		Endpoint:       "localhost:80",
		certsDir:       defaultCertsDir(),
		UserSpaPrefix:  "/u",
		AdminSpaPrefix: "/adm",
		APIPrefix:      "/api",
	}
}

func defaultCertsDir() string {
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
	fs.StringVar(&c.certsDir, "x", c.certsDir,
		`nodeman<->node connection encryption certificates dir.
should contains nodeman.crt nodeman.key ca.crt`)
	fs.StringVar(&c.DBConn, "db", c.DBConn,
		`db connection string`)
	fs.StringVar(&c.UserSpaPrefix, "userspa", c.UserSpaPrefix,
		`user SPA path prefix`)
	fs.StringVar(&c.AdminSpaPrefix, "adminspa", c.AdminSpaPrefix,
		`admin SPA path prefix`)
	fs.StringVar(&c.APIPrefix, "apipref", c.APIPrefix,
		`api path prefix`)
	fs.StringVar(&c.AdminPassword, "admpass", c.AdminPassword,
		`admin password to change`)
	fs.StringVar(&c.JWTSecret, "jwt", c.JWTSecret,
		`jwt secret`)

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

func initCertPaths(c *Config) {
	{
		nodemanCrt := path.Join(c.certsDir, "nodeman.crt")
		if exists, err := checkFileExists(nodemanCrt); err == nil && exists {
			c.nodemanCrt = nodemanCrt
		}
	}
	{
		nodemanKey := path.Join(c.certsDir, "nodeman.key")
		if exists, err := checkFileExists(nodemanKey); err == nil && exists {
			c.nodemanKey = nodemanKey
		}
	}
	{
		rootCrt := path.Join(c.certsDir, "ca.crt")
		if exists, err := checkFileExists(rootCrt); err == nil && exists {
			c.rootCrt = rootCrt
		}
	}
}
