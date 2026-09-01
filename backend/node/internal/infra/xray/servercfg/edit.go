package servercfg

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/tidwall/sjson"
)

func addSrvUsers(cfg string, ins []models.Inbound, us []models.User) (string, error) {
	for _, inbound := range ins {
		cfgUsers := make([]models.CfgUser, 0, len(us))
		for _, u := range us {
			cfgUser, err := inbound.Format.CfgUser(u)
			if err != nil {
				return "", err
			}
			cfgUsers = append(cfgUsers, cfgUser)
		}

		usersPath := fmt.Sprintf("inbounds.#(tag=%s).settings.clients", inbound.Tag)
		withInbound, err := sjson.Set(cfg, usersPath, cfgUsers)
		if err != nil {
			return "", xerr.WrapWithStack(err)
		}
		cfg = withInbound
	}

	return cfg, nil
}

func makeSectionUsers(it models.Inbound, us []models.User) ([]models.CfgUser, error) {
	sectionUsers := make([]models.CfgUser, 0, len(us))
	for _, u := range us {
		cfgUser, err := it.Format.CfgUser(u)
		if err != nil {
			return nil, err
		}
		sectionUsers = append(sectionUsers, cfgUser)
	}
	return sectionUsers, nil
}
