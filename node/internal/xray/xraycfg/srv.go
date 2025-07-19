package xraycfg

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/cfgread"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type ServerCfg struct {
	srvCfg   string
	inbounds []models.Inbound
	apiURL   string
}

func NewServerCfg(path string) (*ServerCfg, error) {
	srvCfg, err := cfgread.ReadJSON(path)
	if err != nil {
		return nil, fmt.Errorf("init srv config: %w", err)
	}

	inbounds := parseSrvInbounds(srvCfg)
	if len(inbounds) == 0 {
		return nil, fmt.Errorf("%w: no supported inbounds in server cfg", errdefs.ErrConfig)
	}

	apiURL := parseSrvApiURL(srvCfg)
	if apiURL == "" {
		return nil, fmt.Errorf("%w: no api url in server cfg", errdefs.ErrConfig)
	}

	return &ServerCfg{
		srvCfg:   srvCfg,
		inbounds: inbounds,
		apiURL:   apiURL,
	}, nil
}

func (cfg *ServerCfg) GetInbounds() []models.Inbound {
	if cfg == nil {
		return []models.Inbound{}
	}
	return cfg.inbounds
}

func (cfg *ServerCfg) GetApiURL() string {
	if cfg == nil {
		//
		return ""
	}
	return cfg.apiURL
}

func (cfg *ServerCfg) GetUsersCfg(users []models.User) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("%w: srv cfg: get users cfg", errdefs.ErrNilObjectCall)
	}

	usersCfg, err := addSrvUsers(cfg.srvCfg, cfg.inbounds, users)
	if err != nil {
		return "", fmt.Errorf("srv cfg: %w", err)
	}
	return usersCfg, nil
}
