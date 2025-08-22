package xraycfg

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/tidwall/sjson"
)

func addSrvUsers(cfg string, ins []models.Inbound, us []models.User) (string, error) {
	usersCfg := cfg
	for _, inbound := range ins {
		sectionUsers, err := makeSectionUsers(inbound.Type, us)
		if err != nil {
			return "", err
		}

		usersPath := fmt.Sprintf("inbounds.#(tag=%s).settings.clients", inbound.Tag)
		usersCfg, err = sjson.Set(usersCfg, usersPath, sectionUsers)
		if err != nil {
			return "", errdefs.WrapWithStack(err)
		}
	}

	return usersCfg, nil
}

func makeSectionUsers(it models.InboundType, us []models.User) ([]map[string]string, error) {
	sectionUsers := make([]map[string]string, 0, len(us))
	for _, u := range us {
		su, err := makeSectionUser(it, u)
		if err != nil {
			return nil, err
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
			"id":    u.VlessUUID,
		}, nil
	case models.VlessXHTTP:
		return map[string]string{
			"email": u.Name,
			"id":    u.VlessUUID,
		}, nil
	default:
		return nil, errdefs.New("unsupported inbound type",
			errdefs.Withf("inbound: %v", it))
	}
}
