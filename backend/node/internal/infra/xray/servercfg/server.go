package servercfg

import (
	"os"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type Config struct {
	config   string
	inbounds []models.Inbound
	apiURL   string
}

func New(path string) (*Config, error) {
	srvCfg, err := os.ReadFile(path)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	srvCfgStr := string(srvCfg)

	inbounds := parseSrvInbounds(srvCfgStr)
	if len(inbounds) == 0 {
		return nil, xerr.New("no supported inbounds in server cfg")
	}

	apiURL := parseSrvApiURL(srvCfgStr)
	if apiURL == "" {
		return nil, xerr.New("no api url in server cfg")
	}

	return &Config{
		config:   srvCfgStr,
		inbounds: inbounds,
		apiURL:   apiURL,
	}, nil
}

func (cfg *Config) GetInbounds() []models.Inbound {
	if cfg == nil {
		return []models.Inbound{}
	}
	return cfg.inbounds
}

func (cfg *Config) GetApiURL() string {
	if cfg == nil {
		return ""
	}
	return cfg.apiURL
}

func (cfg *Config) GetUsersCfg(users []models.User) (string, error) {
	if cfg == nil {
		return "", errdefs.NilCall()
	}

	usersCfg, err := addSrvUsers(cfg.config, cfg.inbounds, users)
	if err != nil {
		return "", err
	}
	return usersCfg, nil
}
