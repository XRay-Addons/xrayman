package config

import (
	"net/url"

	"github.com/XRay-Addons/xrayman/common/xerr"
)

type Config struct {
	Endpoint      string
	DBConn        string
	AdminPassword string
	JwtSecret     string

	ApiServicePath string
	UserSpaPath    string
	AdminSpaPath   string

	ApiServiceUrl string
	UserSpaUrl    string
	AdminSpaUrl   string

	AllowedOrigins []string
}

const apiServicePath = "/api"
const userSpaPath = "/u"
const adminSpaPath = "/adm"

func Init(r RawConfig) (*Config, error) {
	c := Config{
		Endpoint:       r.Endpoint,
		DBConn:         r.DBConn,
		AdminPassword:  r.AdminPassword,
		JwtSecret:      r.JwtSecret,
		ApiServicePath: apiServicePath,
		UserSpaPath:    userSpaPath,
		AdminSpaPath:   adminSpaPath,
	}

	c.ApiServiceUrl = or(r.ApiServiceUrl, c.ApiServicePath)
	c.UserSpaUrl = or(r.UserSpaUrl, c.UserSpaPath)
	c.AdminSpaUrl = or(r.AdminSpaUrl, c.AdminSpaPath)

	for _, u := range []string{r.ApiServiceUrl, r.UserSpaUrl, r.AdminSpaUrl} {
		o, err := getUrlOrigin(u)
		if err != nil {
			return nil, err
		}
		if o != "" {
			c.AllowedOrigins = append(c.AllowedOrigins, o)
		}
	}

	return &c, nil
}

func or(a string, b string) string {
	if a != "" {
		return a
	}
	return b
}

// for empty or relative return ""
// else return origin
func getUrlOrigin(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", xerr.WrapWithStack(err)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return "", nil
	}

	return parsed.Scheme + "://" + parsed.Host, nil
}
