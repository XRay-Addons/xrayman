package xraycfg

import (
	"fmt"
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

func readFile(filePath string) (string, error) {
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("%w: read config file %v", errdefs.ErrConfig, err)
	}
	return string(contentBytes), nil
}
