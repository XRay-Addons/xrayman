package cfgread

import (
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/infra/common/xerr"
	"github.com/tidwall/gjson"
)

func ReadJSON(cfgPath string) (string, error) {
	cfg, err := os.ReadFile(cfgPath) // #nosec
	if err != nil {
		return "", xerr.WrapWithStack(err)
	}
	if !gjson.ValidBytes(cfg) {
		return "", xerr.WrapWithStack(err)
	}
	return string(cfg), nil
}
