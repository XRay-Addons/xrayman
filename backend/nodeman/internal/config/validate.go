package config

import (
	"net"
	"net/url"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
)

func Validate(c *Config) error {
	if err := checkEndpoint(c.Endpoint); err != nil {
		return xerr.InvalidArgf("endpoint: '%s', %v", c.Endpoint, err)
	}
	if err := checkDBConn(c.DBConn); err != nil {
		return xerr.InvalidArgf("db conn: '%s', %v", c.DBConn, err)
	}
	if !checkBaseUrl(c.ApiServiceUrl) {
		return xerr.InvalidArgf("api service url: '%s'", c.ApiServiceUrl)
	}
	if !checkBaseUrl(c.UserSpaUrl) {
		return xerr.InvalidArgf("user spa url: '%s'", c.UserSpaUrl)
	}
	if !checkBaseUrl(c.AdminSpaUrl) {
		return xerr.InvalidArgf("admin spa url: '%s'", c.AdminSpaUrl)
	}
	if !checkSyncInterval(c.StateSyncInterval) {
		return xerr.InvalidArgf("state sync interval: '%d'", c.StateSyncInterval)
	}
	if !checkSyncInterval(c.StateSyncInterval) {
		return xerr.InvalidArgf("stats sync interval: '%d'", c.StatsSyncInterval)
	}
	if !checkJwtSecret(c.JwtSecret) {
		return xerr.InvalidArgf("jwt secret: '%s'", c.JwtSecret)
	}
	if c.MetricsEndpoint != "" {
		if err := checkEndpoint(c.MetricsEndpoint); err != nil {
			return xerr.InvalidArgf("metrics endpoint: '%s', %v", c.MetricsEndpoint, err)
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
