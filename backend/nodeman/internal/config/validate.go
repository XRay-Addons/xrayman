package config

import (
	"net"
	"strings"

	"github.com/XRay-Addons/xrayman/common/xerr"
)

func Validate(c Config) error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return xerr.New("invalid endpoint", xerr.Withf("endpoint: %s", c.Endpoint))
	}
	if err := checkDBConn(c); err != nil {
		return err
	}
	if err := checkPaths(c); err != nil {
		return err
	}
	if err := checkAuth(c); err != nil {
		return err
	}

	return nil
}

func checkDBConn(c Config) error {
	if len(c.DBConn) == 0 {
		return xerr.New("dbconn string invalid")
	}
	return nil
}

func checkPaths(c Config) error {
	if !strings.HasPrefix(c.APIPrefix, "/") {
		return xerr.New("api prefix invalid")
	}
	if !strings.HasPrefix(c.UserSpaPrefix, "/") {
		return xerr.New("user spa prefix invalid")
	}
	if !strings.HasPrefix(c.AdminSpaPrefix, "/") {
		return xerr.New("admin spa prefix invalid")
	}
	return nil
}

func checkAuth(c Config) error {
	if c.JWTSecret == "" {
		return xerr.New("jwt secret invalid")
	}
	return nil
}
