package xraycfg

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/tidwall/gjson"
)

const (
	inboundsPath   = "inbounds"
	protocolPath   = "protocol"
	inboundTagPath = "tag"
	networkPath    = "streamSettings.network"
	securityPath   = "streamSettings.security"
)

func GetInbounds(serverConfig string) ([]models.Inbound, error) {
	if !gjson.Valid(serverConfig) {
		return nil, fmt.Errorf("%w: invalid server config json", errdefs.ErrConfig)
	}

	inboundSections := gjson.Get(serverConfig, inboundsPath).Array()

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

	return inbounds, nil
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
