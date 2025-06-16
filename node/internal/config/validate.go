package config

import (
	"fmt"
	"net"
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

func Validate(c Config) error {
	if _, err := os.Stat(c.XRayExecPath); err != nil {
		return fmt.Errorf("%w: xray exec path: %v", errdefs.ErrConfig, err)
	}
	if _, err := os.Stat(c.XRayServerConfigPath); err != nil {
		return fmt.Errorf("%w: server config path: %v", errdefs.ErrConfig, err)
	}
	if _, err := os.Stat(c.XRayClientConfigPath); err != nil {
		return fmt.Errorf("%w: client config path: %v", errdefs.ErrConfig, err)
	}
	if _, err := net.ResolveTCPAddr("tcp", c.Endpoint); err != nil {
		return fmt.Errorf("%w: invalid endpoint: %v", err)
	}
	return nil
}
