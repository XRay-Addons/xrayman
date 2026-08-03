package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/caarlos0/env/v6"
	"github.com/kr/text"
	"go.uber.org/zap/zapcore"
)

func LoadConfig() (*RawConfig, error) {
	cfg := defaultConfig()
	if err := readCLIParams(cfg); err != nil {
		return nil, err
	}
	if err := readEnvParams(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func defaultConfig() *RawConfig {
	return &RawConfig{
		Endpoint:           "localhost:80",
		StateSyncInterval:  5,
		StatsSyncInterval:  60,
		StorageCallTimeout: 5,
		NodeCallTimeout:    5,
		LogLevel:           zapcore.InfoLevel.String(),
	}
}

func readCLIParams(c *RawConfig) error {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	fs.StringVar(&c.Endpoint, "a", c.Endpoint,
		"server endpoint tcp address, like :8080, 127.0.0.1:80, localhost:22")
	fs.StringVar(&c.DBConn, "db", c.DBConn,
		"postgress connection string, like postgresql://user@password/127.0.0.1:4321/dbname")
	fs.StringVar(&c.JwtSecret, "jwt", c.JwtSecret,
		"jwt secret")

	fs.IntVar(&c.StateSyncInterval, "state", c.StateSyncInterval,
		"state sync interval, s")
	fs.IntVar(&c.StatsSyncInterval, "stats", c.StatsSyncInterval,
		"stats sync interval, s")

	fs.StringVar(&c.ApiServiceUrl, "apisrv", c.ApiServiceUrl,
		`public base URL of the API as seen by browsers (used for CORS and SPAs config).
If empty or relative, the internal API base path is used.
should be like /internal/api or https://api.example.com/api (optional)`)
	fs.StringVar(&c.UserSpaUrl, "userspa", c.UserSpaUrl,
		`public base URL of the User SPA as seen by browsers (used for CORS and SPAs config).
If empty or relative, the internal User SPA base path is used.
should be like /user or https://u.example.com (optional)`)
	fs.StringVar(&c.AdminSpaUrl, "adminspa", c.AdminSpaUrl,
		`public base URL of the Admin SPA as seen by browsers (used for CORS and SPAs config).
If empty or relative, the internal Admin SPA base path is used.
should be like /admin or https://adm.example.com (optional)`)

	fs.StringVar(&c.AdminPassword, "admpass", c.AdminPassword,
		`admin password to change (optional)`)

	fs.IntVar(&c.StorageCallTimeout, "storage-timeout", c.StorageCallTimeout,
		`storage call timeout, s (optional)`)
	fs.IntVar(&c.NodeCallTimeout, "node-timeout", c.NodeCallTimeout,
		`node call timeout, s (optional)`)

	fs.StringVar(&c.LogLevel, "log-lvl", c.LogLevel,
		`zap log level (optional)`)

	fs.Usage = func() {
		fmt.Printf("Usage:\n")
		argGroups := [][]string{
			{"a", "db", "jwt"},
			{"state", "stats"},
			{"apisrv", "userspa", "adminspa"},
			{"admpass"},
			{"storage-timeout", "node-timeout"},
			{"log-lvl"},
		}

		for _, argGroup := range argGroups {
			for _, arg := range argGroup {
				flag := fs.Lookup(arg)
				fmt.Printf(" -%s\n", flag.Name)
				fmt.Printf("%s\n", text.Indent(flag.Usage, "    "))
				if len(flag.DefValue) > 0 {
					fmt.Printf("    default: %s\n", flag.DefValue)
				}
			}
			fmt.Printf("\n")
		}
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		return xerr.WrapWithStack(err)
	}

	return nil
}

func readEnvParams(c *RawConfig) error {
	if err := env.Parse(c); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}
