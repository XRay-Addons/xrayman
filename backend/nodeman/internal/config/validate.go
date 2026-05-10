package config

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
)

func Validate(c Config) error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return xerr.New("invalid endpoint", xerr.Withf("endpoint: %s", c.Endpoint))
	}
	if err := checkCerts(c); err != nil {
		return err
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

func checkFileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	if !info.Mode().IsRegular() {
		return false, xerr.New("file is not regular",
			errdefs.WithFile(path))
	}
	return true, nil
}

func checkCerts(c Config) error {
	// all 3 certs should exists or not together
	certs := []string{c.nodemanCrt, c.nodemanKey, c.rootCrt}
	existsCount := 0
	existsDescription := ""
	for _, f := range []string{c.nodemanCrt, c.nodemanKey, c.rootCrt} {
		exists, err := checkFileExists(f)
		if err != nil {
			return err
		}
		existsDescription += fmt.Sprintf("cert %s: %v", f, exists)
		if exists {
			existsCount++
		}
	}

	if existsCount == 0 || existsCount == len(certs) {
		return nil
	}

	return xerr.New("cert files inconsistency",
		xerr.Withf("inconsistency: %s", existsDescription))
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
