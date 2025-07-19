package xraycfg

import (
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/tidwall/gjson"
)

const (
	inboundsPath   = "inbounds"
	protocolPath   = "protocol"
	inboundTagPath = "tag"
	networkPath    = "streamSettings.network"
	securityPath   = "streamSettings.security"
	apiUrlPath     = "api.listen"
)

func parseSrvInbounds(cfg string) []models.Inbound {
	inboundSections := gjson.Get(cfg, inboundsPath).Array()

	inbounds := make([]models.Inbound, 0, len(inboundSections))
	for _, inbound := range inboundSections {
		tag := inbound.Get(inboundTagPath).String()
		protocol := inbound.Get(protocolPath).String()
		network := inbound.Get(networkPath).String()
		security := inbound.Get(securityPath).String()

		inboundType := getInboundType(protocol, network, security)
		if inboundType == models.UnsupportedInbound {
			continue
		}
		inbounds = append(inbounds, models.Inbound{
			Tag:  tag,
			Type: inboundType,
		})
	}

	return inbounds
}

func getInboundType(protocol, network, security string) models.InboundType {
	if protocol != "vless" {
		return models.UnsupportedInbound
	}
	if network == "tcp" && security == "reality" {
		return models.VlessTcpReality
	}
	if network == "xhttp" {
		return models.VlessXHTTP
	}
	return models.UnsupportedInbound
}

func parseSrvApiURL(srvCfg string) string {
	apiURL := gjson.Get(srvCfg, apiUrlPath).String()
	return apiURL
}
