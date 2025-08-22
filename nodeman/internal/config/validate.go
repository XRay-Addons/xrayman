package config

import (
	"fmt"
	"net"
	"os"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
)

func Validate(c Config) error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return errdefs.New("invalid endpoint", errdefs.Withf("endpoint: %s", c.Endpoint))
	}
	if err := checkCerts(c); err != nil {
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
		return false, errdefs.New("file is not regular",
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

	return errdefs.New("cert files inconsistency",
		errdefs.Withf("inconsistency: %s", existsDescription))
}
