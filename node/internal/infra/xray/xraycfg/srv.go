package xraycfg

import (
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/infra/cfgread"
	"github.com/XRay-Addons/xrayman/node/internal/infra/common/xerr"
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
		return nil, err
	}

	inbounds := parseSrvInbounds(srvCfg)
	if len(inbounds) == 0 {
		return nil, xerr.New("no supported inbounds in server cfg")
	}

	apiURL := parseSrvApiURL(srvCfg)
	if apiURL == "" {
		return nil, xerr.New("no api url in server cfg")
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
		return ""
	}
	return cfg.apiURL
}

func (cfg *ServerCfg) GetUsersCfg(users []models.User) (string, error) {
	if cfg == nil {
		return "", errdefs.NilCall()
	}

	usersCfg, err := addSrvUsers(cfg.srvCfg, cfg.inbounds, users)
	if err != nil {
		return "", err
	}
	return usersCfg, nil
}
