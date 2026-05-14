package config

import (
	"flag"
	"os"

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
	return &Config{
		Endpoint:       "localhost:80",
		UserSpaPrefix:  "/u",
		AdminSpaPrefix: "/adm",
		APIPrefix:      "/api",
	}
}

func readCLIParams(c *Config) error {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	fs.StringVar(&c.Endpoint, "a", c.Endpoint,
		"server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22")
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
