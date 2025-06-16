package xraycfg

import (
	"fmt"
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type XRayCfg struct {
	serverCfg string
	clientCfg string
	inbounds  []models.Inbound
	apiURL    string
}

func New(serverCfgPath, clientCfgPath string) (*XRayCfg, error) {
	serverCfg, err := readFile(serverCfgPath)
	if err != nil {
		return nil, fmt.Errorf("%w: server config file reading: %v", errdefs.ErrConfig, err)
	}

	clientCfg, err := readFile(clientCfgPath)
	if err != nil {
		return nil, fmt.Errorf("%w: user config file reading: %v", errdefs.ErrConfig, err)
	}

	inbounds, err := GetInbounds(serverCfg)
	if err != nil {
		return nil, fmt.Errorf("get cfg inbounds: %w", err)
	}
	apiURL, err := GetApiURL(serverCfg)
	if err != nil {
		return nil, fmt.Errorf("get api url: %w", err)
	}

	return &XRayCfg{
		serverCfg: serverCfg,
		clientCfg: clientCfg,
		inbounds:  inbounds,
		apiURL:    apiURL,
	}, nil
}

func (cfg *XRayCfg) GetInbounds() []models.Inbound {
	if cfg == nil {
		return nil
	}
	return cfg.inbounds
}

func (cfg *XRayCfg) GetApiURL() string {
	if cfg == nil {
		return ""
	}
	return cfg.apiURL
}

func (cfg *XRayCfg) GetServerConfig(users []models.User) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("%w: xray cfg not exists", errdefs.ErrIPE)
	}

	usersConfig, err := AddUsers(cfg.serverCfg, cfg.inbounds, users)
	if err != nil {
		return "", fmt.Errorf("add config users: %w", err)
	}
	return usersConfig, nil
}

func (cfg *XRayCfg) GetClientConfig() string {
	if cfg == nil {
		return ""
	}
	return cfg.clientCfg
}

func readFile(filePath string) (string, error) {
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read file %s: %w", filePath, err)
	}
	return string(contentBytes), nil
}
