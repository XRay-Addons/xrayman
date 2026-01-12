package cfgread

import (
	"fmt"
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/tidwall/gjson"
)

func ReadJSON(cfgPath string) (string, error) {
	cfg, err := os.ReadFile(cfgPath)
	if err != nil {
		return "", fmt.Errorf("%w: read config file %v", errdefs.ErrConfig, err)
	}
	if !gjson.ValidBytes(cfg) {
		return "", fmt.Errorf("%w: invalid config json", errdefs.ErrConfig)
	}
	return string(cfg), nil
}
