package models

import "github.com/xtls/xray-core/common/protocol"

type CfgUser = map[string]string

// inbound format is combination of protocol, network and security
type InboundFormat interface {
	Check(protocol, network, security string) bool
	CfgUser(u User) (CfgUser, error)
	ApiUser(u User) (*protocol.User, error)
}

type Inbound struct {
	Tag    string
	Format InboundFormat
}
