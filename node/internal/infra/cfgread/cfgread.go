package cfgread

import (
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/tidwall/gjson"
)

func ReadJSON(cfgPath string) (string, error) {
	cfg, err := os.ReadFile(cfgPath)
	if err != nil {
		return "", errdefs.WithStack(err)
	}
	if !gjson.ValidBytes(cfg) {
		return "", errdefs.WithStack(err)
	}
	return string(cfg), nil
}
