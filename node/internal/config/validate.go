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
	if len(c.AccessKey) == 0 {
		return fmt.Errorf("%w: invalid access key: %v", errdefs.ErrConfig, c.AccessKey)
	}
	if err := checkExecutable(c.XRayExecPath); err != nil {
		return err
	}
	if err := checkFile(c.XRayConfigPath); err != nil {
		return err
	}
	return nil
}

func checkExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%w: xray exec file not exists", errdefs.ErrConfig)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%w: xray exec is not regular file", errdefs.ErrConfig)
	}
	if info.Mode().Perm()&0111 != 0 {
		return fmt.Errorf("%w: xray exec is not executable for current user", errdefs.ErrConfig)
	}
	return nil
}

func checkFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%w: xray exec file not exists", errdefs.ErrConfig)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%w: xray exec is not regular file", errdefs.ErrConfig)
	}
	return nil
}
