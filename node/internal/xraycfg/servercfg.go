package xraycfg

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/tidwall/gjson"
)

type ServerCfg struct {
	serverCfg string
	inbounds  []models.Inbound
	apiURL    string
}

func NewServerCfg(serverCfgPath string) (*ServerCfg, error) {
	serverCfg, err := readFile(serverCfgPath)
	if err != nil {
		return nil, fmt.Errorf("%w: server config file reading: %v", errdefs.ErrConfig, err)
	}
	if !gjson.Valid(serverCfg) {
		return nil, fmt.Errorf("%w: invalid server config json", errdefs.ErrConfig)
	}

	inbounds := parseServerInbounds(serverCfg)
	if len(inbounds) == 0 {
		return nil, fmt.Errorf("%w: no supported inbounds in server cfg", errdefs.ErrConfig)
	}

	apiURL := parseServerApiURL(serverCfg)
	if apiURL == "" {
		return nil, fmt.Errorf("%w: no api url in server cfg", errdefs.ErrConfig)
	}

	return &ServerCfg{
		serverCfg: serverCfg,
		inbounds:  inbounds,
		apiURL:    apiURL,
	}, nil
}

func (cfg *ServerCfg) GetInbounds() ([]models.Inbound, error) {
	if cfg == nil {
		return nil, fmt.Errorf("%w: server cfg", errdefs.ErrNilObjectCall)
	}
	return cfg.inbounds, nil
}

func (cfg *ServerCfg) GetApiURL() (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("%w: server cfg", errdefs.ErrNilObjectCall)
	}
	return cfg.apiURL, nil
}

func (cfg *ServerCfg) GetServerConfig(users []models.User) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("%w: xray cfg not exists", errdefs.ErrIPE)
	}

	usersConfig, err := addServerConfigUsers(cfg.serverCfg, cfg.inbounds, users)
	if err != nil {
		return "", fmt.Errorf("get server config: %w", err)
	}
	return usersConfig, nil
}
