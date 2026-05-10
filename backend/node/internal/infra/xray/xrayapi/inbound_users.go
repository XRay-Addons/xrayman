package xrayapi

import (
	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/vless"
)

func getInboundUser(u models.User, in models.InboundType) (*protocol.User, error) {
	switch in {
	case models.VlessTcpReality, models.VlessXHTTP:
		vlessAccunt, err := getVlessAccount(u, in)
		if err != nil {
			return nil, err
		}
		return &protocol.User{
			Email:   u.VlessEmail(),
			Account: serial.ToTypedMessage(vlessAccunt),
		}, nil
	default:
		return nil, xerr.Newf("unsupported inbound: %v", in)
	}
}

func getVlessAccount(u models.User, in models.InboundType) (*vless.Account, error) {
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
		return nil, xerr.Newf("unsupported inbound: %v", in)
	}
}
