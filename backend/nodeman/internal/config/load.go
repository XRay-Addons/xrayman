package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/caarlos0/env/v6"
	"github.com/kr/text"
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
		Endpoint: "localhost:80",
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

	fs.Usage = func() {
		fmt.Printf("Usage:\n")
		argGroups := [][]string{
			{"a", "db", "jwt"},
			{"apisrv", "userspa", "adminspa"},
			{"admpass"},
		}

		for _, argGroup := range argGroups {
			for _, arg := range argGroup {
				flag := fs.Lookup(arg)
				fmt.Printf(" -%s\n", flag.Name)
				fmt.Printf("%s\n", text.Indent(flag.Usage, "    "))
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
