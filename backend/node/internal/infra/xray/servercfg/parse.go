package servercfg

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
)

const (
	inboundsPath   = "inbounds"
	protocolPath   = "protocol"
	inboundTagPath = "tag"
	networkPath    = "streamSettings.network"
	securityPath   = "streamSettings.security"
	apiUrlPath     = "api.listen"
)

func parseSrvInbounds(cfg string, fmts []models.InboundFormat,
	log *zap.Logger,
) []models.Inbound {
	inboundSections := gjson.Get(cfg, inboundsPath).Array()

	inbounds := make([]models.Inbound, 0, len(inboundSections))
	for _, inbound := range inboundSections {
		tag := inbound.Get(inboundTagPath).String()
		protocol := inbound.Get(protocolPath).String()
		network := inbound.Get(networkPath).String()
		security := inbound.Get(securityPath).String()

		if f, ok := findInboundFormat(fmts, protocol, network, security); ok {
			inbounds = append(inbounds, models.Inbound{Tag: tag, Format: f})
		} else {
			log.Warn(fmt.Sprintf(
				"unsupported inbound format '%s': %v(protocol)-%v(network)-%v(security)",
				tag, protocol, network, security))
		}
	}

	return inbounds
}

func findInboundFormat(fmts []models.InboundFormat, p, n, s string) (models.InboundFormat, bool) {
	for _, f := range fmts {
		if f.Check(p, n, s) {
			return f, true
		}
	}
	return nil, false
}

func parseSrvApiURL(srvCfg string) string {
	apiURL := gjson.Get(srvCfg, apiUrlPath).String()
	return apiURL
}
