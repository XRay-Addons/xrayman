package config

import (
	"net"
	"net/url"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
)

func Validate(c *Config) error {
	if err := checkEndpoint(c.Endpoint); err != nil {
		return xerr.InvalidArg("endpoint", c.Endpoint, err)
	}
	if err := checkDBConn(c.DBConn); err != nil {
		return xerr.InvalidArg("db conn", c.DBConn, err)
	}
	if !checkBaseUrl(c.ApiServiceUrl) {
		return xerr.InvalidArg("api service url", c.ApiServiceUrl, nil)
	}
	if !checkBaseUrl(c.UserSpaUrl) {
		return xerr.InvalidArg("user spa url", c.UserSpaUrl, nil)
	}
	if !checkBaseUrl(c.AdminSpaUrl) {
		return xerr.InvalidArg("admin spa url", c.AdminSpaUrl, nil)
	}
	if !checkSyncInterval(c.StateSyncInterval) {
		return xerr.InvalidArg("state sync interval", c.StateSyncInterval, nil)
	}
	if !checkSyncInterval(c.StateSyncInterval) {
		return xerr.InvalidArg("stats sync interval", c.StatsSyncInterval, nil)
	}
	if !checkJwtSecret(c.JwtSecret) {
		return xerr.InvalidArg("jwt secret", c.JwtSecret, nil)
	}
	if c.MetricsEndpoint != "" {
		if err := checkEndpoint(c.MetricsEndpoint); err != nil {
			return xerr.InvalidArg("metrics endpoint", c.MetricsEndpoint, err)
		}
	}

	return nil
}

func checkEndpoint(e string) error {
	_, err := net.ResolveTCPAddr("tcp", e)
	return err
}

func checkDBConn(dbconn string) error {
	if len(dbconn) == 0 {
		return xerr.New("dbconn string invalid")
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

func checkSyncInterval(interval time.Duration) bool {
	return interval > 0
}

func checkJwtSecret(s string) bool {
	return s != ""
}
