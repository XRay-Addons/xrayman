package xray

import (
	"fmt"
	"os"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/tidwall/gjson"
)

type InboundType int

const (
	UnsupportedInbound InboundType = iota
	VlessTcpReality                = iota
	VlessXHTTP
)

type InboundUsersSection struct {
	JSONPath string
	Type     InboundType
}

func CreateServerConfig(emptyConfigPath string, users []User) (string, error) {
	emptyConfig, err := readTextFile(emptyConfigPath)
	if err != nil {
		return "", fmt.Errorf("%w: read server config %v", errdefs.ErrConfig, err)
	}
	inboundUserSections, err := parseUserSections(emptyConfig)
	if err != nil {
		return "", fmt.Errorf("%w: parse inbound sections %v", errdefs.ErrConfig, err)
	}
	filledConfig, err := insertConfigUsers(emptyConfig, inboundUserSections, users)
}

func readTextFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

const (
	inboundsPath = "inbounds"
	protocolPath = "protocol"
	networkPath  = "streamSettings.network"
	securityPath = "streamSettings.security"
)

func parseUserSections(config string) ([]InboundUsersSection, error) {
	inbounds := gjson.Get(config, inboundsPath).Array()
	if len(inbounds) == 0 {
		return nil, fmt.Errorf("server config contains no inbounds")
	}

	inboundSections := make([]InboundUsersSection, 0, len(inbounds))
	for idx, inbound := range inbounds {
		protocol := inbound.Get(protocolPath).String()
		network := inbound.Get(networkPath).String()
		security := inbound.Get(securityPath).String()

		inboundType := getInboundType(protocol, network, security)
		if inboundType == UnsupportedInbound {
			continue
		}

		inboundSection := InboundUsersSection{
			JSONPath: fmt.Sprintf("%s.#d.settings.clients", inboundsPath, idx),
			Type:     inboundType,
		}
		inboundSections = append(inboundSections, inboundSection)
	}

	return inboundSections, nil
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

func insertConfigUsers(emptyConfig string, inboundSections []InboundUsersSection, users []User) (string, error) {
	return "", nil
}
