package formats

import (
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/proxy/vless"
)

type VlessXHTTP struct {
}

var _ models.InboundFormat = (*VlessXHTTP)(nil)

func (v *VlessXHTTP) Check(protocol, network, security string) bool {
	return protocol == "vless" && network == "xhttp"
}

func (v *VlessXHTTP) CfgUser(u models.User) (models.CfgUser, error) {
	return map[string]string{
		"email": u.VlessEmail(),
		"id":    u.VlessUUID,
	}, nil
}

func (v *VlessXHTTP) ApiUser(u models.User) (*protocol.User, error) {
	return &protocol.User{
		Email: u.VlessEmail(),
		Account: serial.ToTypedMessage(&vless.Account{
			Id:         u.VlessUUID,
			Encryption: "none",
		}),
	}, nil
}
