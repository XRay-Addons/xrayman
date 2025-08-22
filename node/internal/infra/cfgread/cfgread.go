package cfgread

import (
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/tidwall/gjson"
)

func ReadJSON(cfgPath string) (string, error) {
	cfg, err := os.ReadFile(cfgPath) // #nosec
	if err != nil {
		return "", errdefs.WrapWithStack(err)
	}
	if !gjson.ValidBytes(cfg) {
		return "", errdefs.WrapWithStack(err)
	}
	return string(cfg), nil
}
