package models

type NodeID = int

type NodeClientCfg struct {
	Template       string
	UserNameField  string
	VlessUUIDField string
}

type NodeProperties struct {
	ClientCfg NodeClientCfg
}

type Node struct {
	ID         NodeID
	Endpoint   string
	Properties NodeProperties
}

type NodeStatus int

const (
	NodeStatusUnknown NodeStatus = iota + 1
	NodeOff
	NodeOn
)

func (s NodeStatus) String() string {
	switch s {
	case NodeOff:
		return "Off"
	case NodeOn:
		return "On"
	default:
		return "Unknown"
	}
}
