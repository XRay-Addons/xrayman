package formats

import (
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/vless"
)

type VlessTCPReality struct {
}

var _ models.InboundFormat = (*VlessTCPReality)(nil)

func (v *VlessTCPReality) Check(protocol, network, security string) bool {
	return protocol == "vless" && network == "tcp" && security == "reality"
}

func (v *VlessTCPReality) CfgUser(u models.User) (models.CfgUser, error) {
	return map[string]string{
		"email": u.VlessEmail(),
		"flow":  "xtls-rprx-vision",
		"id":    u.VlessUUID,
	}, nil
}

func (v *VlessTCPReality) ApiUser(u models.User) (*protocol.User, error) {
	return &protocol.User{
		Email: u.VlessEmail(),
		Account: serial.ToTypedMessage(&vless.Account{
			Id:         u.VlessUUID,
			Encryption: "none",
			Flow:       "xtls-rprx-vision",
		}),
	}, nil
}
