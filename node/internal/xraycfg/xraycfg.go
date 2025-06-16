package xraycfg

import (
	"fmt"
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type XRayCfg struct {
	cfgPath    string
	cfgContent string
	inbounds   []models.Inbound
	apiURL     string
}

func New(cfgPath string) (*XRayCfg, error) {
	cfgContentBytes, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("%w: config file reading: %v", errdefs.ErrConfig, err)
	}
	cfgContent := string(cfgContentBytes)

	inbounds, err := GetInbounds(cfgContent)
	if err != nil {
		return nil, fmt.Errorf("get cfg inbounds: %w", err)
	}
	apiURL, err := GetApiURL(cfgContent)
	if err != nil {
		return nil, fmt.Errorf("get api url: %w", err)
	}

	return &XRayCfg{
		cfgPath:    cfgPath,
		cfgContent: cfgContent,
		inbounds:   inbounds,
		apiURL:     apiURL,
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

func (cfg *XRayCfg) SetUsers(users []models.User) error {
	if cfg == nil {
		return fmt.Errorf("%w: xray cfg not exists", errdefs.ErrIPE)
	}

	editedCfg, err := AddUsers(cfg.cfgContent, cfg.inbounds, users)
	if err != nil {
		return fmt.Errorf("add config users: %w", err)
	}
	if err := os.WriteFile(cfg.cfgPath, []byte(editedCfg), 0o644); err != nil {
		return fmt.Errorf("%w: config file write: %v", errdefs.ErrAccess, err)
	}
	return nil
}
