package xraycfg

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func AddUsers(serverConfig string, inbounds []models.Inbound, users []models.User) (string, error) {
	if !gjson.Valid(serverConfig) {
		return "", fmt.Errorf("%w: invalid server config json", errdefs.ErrConfig)
	}
	filledConfig, err := insertUsers(serverConfig, inbounds, users)
	if err != nil {
		return "", fmt.Errorf("%w: insert config users: %v", errdefs.ErrConfig, err)
	}
	return filledConfig, nil
}

func insertUsers(emptyConfig string, inbounds []models.Inbound, users []models.User) (string, error) {
	config := emptyConfig
	for _, inbound := range inbounds {
		sectionUsers, err := makeSectionUsers(inbound.Type, users)
		if err != nil {
			return "", fmt.Errorf("make section user: %w", err)
		}

		usersPath := fmt.Sprintf("inbounds.#(tag=%s).settings.clients", inbound.Tag)
		config, err = sjson.Set(config, usersPath, sectionUsers)
		if err != nil {
			return "", fmt.Errorf("set config users: %w", err)
		}
	}

	return config, nil
}

func makeSectionUsers(it models.InboundType, users []models.User) ([]map[string]string, error) {
	sectionUsers := make([]map[string]string, 0, len(users))
	for _, u := range users {
		su, err := makeSectionUser(it, u)
		if err != nil {
			return nil, fmt.Errorf("make section user: %w", err)
		}
		sectionUsers = append(sectionUsers, su)
	}
	return sectionUsers, nil
}

func makeSectionUser(it models.InboundType, u models.User) (map[string]string, error) {
	switch it {
	case models.VlessTcpReality:
		return map[string]string{
			"email": u.Name,
			"flow":  "xtls-rprx-vision",
			"id":    u.UUID,
		}, nil
	case models.VlessXHTTP:
		return map[string]string{
			"email": u.Name,
			"id":    u.UUID,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported inbound type")
	}
}
