package xraycfg

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/tidwall/gjson"
)

const (
	apiUrlPath = "api.listen"
)

func GetApiURL(serverConfig string) (string, error) {
	if !gjson.Valid(serverConfig) {
		return "", fmt.Errorf("%w: invalid server config json", errdefs.ErrConfig)
	}

	apiURL := gjson.Get(serverConfig, apiUrlPath).String()
	if apiURL == "" {
		return "", fmt.Errorf("%w: empty api url", errdefs.ErrConfig)
	}
	return apiURL, nil
}
