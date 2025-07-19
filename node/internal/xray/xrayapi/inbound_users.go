package xrayapi

import (
	"fmt"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/vless"
)

func getInboundUser(u models.User, in models.InboundType) (*protocol.User, error) {
	switch in {
	case models.VlessTcpReality, models.VlessXHTTP:
		vlessAccunt, err := getVlessAccound(u, in)
		if err != nil {
			return nil, err
		}
		return &protocol.User{
			Email:   u.Name,
			Account: serial.ToTypedMessage(vlessAccunt),
		}, nil
	default:
		return nil, fmt.Errorf("%w: unsupported inbound: %v", errdefs.ErrConfig, in)
	}
}

func getVlessAccound(u models.User, in models.InboundType) (*vless.Account, error) {
	switch in {
	case models.VlessTcpReality:
		return &vless.Account{
			Id:         u.VlessUUID,
			Encryption: "none",
			Flow:       "xtls-rprx-vision",
		}, nil
	case models.VlessXHTTP:
		return &vless.Account{
			Id:         u.VlessUUID,
			Encryption: "none",
		}, nil
	default:
		return nil, fmt.Errorf("%w: unsupported inbound", errdefs.ErrIPE)
	}
}
