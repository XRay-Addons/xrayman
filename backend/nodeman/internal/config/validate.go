package config

import (
	"net"
	"net/url"

	"github.com/XRay-Addons/xrayman/common/xerr"
)

func Validate(c RawConfig) error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return xerr.New("invalid endpoint", xerr.Withf("endpoint: %s", c.Endpoint))
	}
	if err := checkDBConn(c); err != nil {
		return err
	}
	if err := checkBaseUrls(c); err != nil {
		return err
	}
	if err := checkAuth(c); err != nil {
		return err
	}

	return nil
}

func checkDBConn(c RawConfig) error {
	if len(c.DBConn) == 0 {
		return xerr.New("dbconn string invalid")
	}
	return nil
}

func checkBaseUrls(c RawConfig) error {
	if !checkBaseUrl(c.ApiServiceUrl) {
		return xerr.New("api service invalid")
	}
	if !checkBaseUrl(c.UserSpaUrl) {
		return xerr.New("user spa url invalid")
	}
	if !checkBaseUrl(c.AdminSpaUrl) {
		return xerr.New("admin spa url invalid")
	}
	return nil
}

// check if u = schema://host/path or /path or empty
func checkBaseUrl(u string) bool {
	parsed, err := url.Parse(u)
	if err != nil {
		return false
	}
	return (parsed.Scheme == "" && parsed.Host == "") ||
		(parsed.Scheme != "" && parsed.Host != "")
}

func checkAuth(c RawConfig) error {
	if c.JwtSecret == "" {
		return xerr.New("jwt secret invalid")
	}
	return nil
}
