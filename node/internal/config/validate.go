package config

import (
	"fmt"
	"net"
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

func Validate(c Config) error {
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return fmt.Errorf("%w: invalid endpoint: %v", errdefs.ErrConfig, err)
	}
	if len(c.AccessKey) != 0 && len(c.AccessKey) != 32 {
		return fmt.Errorf("%w: invalid access key length %v, required 32",
			errdefs.ErrConfig, c.AccessKey)
	}
	if err := checkExecutable(c.XRayExec()); err != nil {
		return fmt.Errorf("%w: xray exec: %v", errdefs.ErrConfig, err)
	}
	if err := checkFile(c.XRayServer()); err != nil {
		return fmt.Errorf("%w: xray server cfg: %v", errdefs.ErrConfig, err)
	}
	if err := checkFile(c.XRayClient()); err != nil {
		return fmt.Errorf("%w: xray client cfg: %v", errdefs.ErrConfig, err)
	}

	return nil
}

func checkExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s file not exists", path)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s not a regular file", path)
	}
	perm := info.Mode().Perm()
	if perm&0111 != 0 {
		return nil
	}
	return fmt.Errorf("%s file not executable for current user: %v", path, perm)
}

func checkFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s file not exists", path)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", path)
	}
	return nil
}
