package xraycfg

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type InboundType int

const (
	UnsupportedInbound InboundType = iota
	VlessTcpReality                = iota
	VlessXHTTP
)

type UsersSection struct {
	JSONPath string
	Type     InboundType
}

func CreateServerConfig(emptyConfig string, users []models.User) (string, error) {
	if !gjson.Valid(emptyConfig) {
		return "", fmt.Errorf("%w: invalid server config json", errdefs.ErrConfig)
	}
	inboundUserSections, err := parseUserSections(emptyConfig)
	if err != nil {
		return "", fmt.Errorf("%w: parse inbound sections: %v", errdefs.ErrConfig, err)
	}
	filledConfig, err := insertConfigUsers(emptyConfig, inboundUserSections, users)
	if err != nil {
		return "", fmt.Errorf("%w: insert config users: %v", errdefs.ErrConfig, err)
	}
	return filledConfig, nil
}

const (
	inboundsPath = "inbounds"
	protocolPath = "protocol"
	networkPath  = "streamSettings.network"
	securityPath = "streamSettings.security"
)

func parseUserSections(config string) ([]UsersSection, error) {
	inbounds := gjson.Get(config, inboundsPath).Array()
	if len(inbounds) == 0 {
		return nil, fmt.Errorf("server config contains no inbounds")
	}

	userSections := make([]UsersSection, 0, len(inbounds))
	for idx, inbound := range inbounds {
		protocol := inbound.Get(protocolPath).String()
		network := inbound.Get(networkPath).String()
		security := inbound.Get(securityPath).String()

		inboundType := getInboundType(protocol, network, security)
		if inboundType == UnsupportedInbound {
			continue
		}

		userSections = append(userSections, UsersSection{
			JSONPath: fmt.Sprintf("%s.%d.settings.clients", inboundsPath, idx),
			Type:     inboundType,
		})
	}

	return userSections, nil
}

func getInboundType(protocol, network, security string) InboundType {
	if protocol != "vless" {
		return UnsupportedInbound
	}
	if network == "tcp" && security == "reality" {
		return VlessTcpReality
	}
	if network == "xhttp" {
		return VlessXHTTP
	}
	return UnsupportedInbound
}

func insertConfigUsers(emptyConfig string, sections []UsersSection, users []models.User) (string, error) {
	config := emptyConfig
	for _, s := range sections {
		sectionUsers, err := makeSectionUsers(s.Type, users)
		if err != nil {
			return "", fmt.Errorf("make section user: %w", err)
		}
		config, err = sjson.Set(config, s.JSONPath, sectionUsers)
		if err != nil {
			return "", fmt.Errorf("set config users: %w", err)
		}
	}

	return config, nil
}

func makeSectionUsers(it InboundType, users []models.User) ([]map[string]string, error) {
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

func makeSectionUser(it InboundType, u models.User) (map[string]string, error) {
	switch it {
	case VlessTcpReality:
		return map[string]string{
			"email": u.Name,
			"flow":  "xtls-rprx-vision",
			"id":    u.UUID,
		}, nil
	case VlessXHTTP:
		return map[string]string{
			"email": u.Name,
			"id":    u.UUID,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported inbound type")
	}
}
