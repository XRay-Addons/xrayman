package models

type NodeStatus int

const (
	NodeStatusUnknown NodeStatus = iota + 1
	NodeStatusStopped
	NodeStatusRunning
)

type ClientConfig struct {
	Template       string
	UserNameField  string
	VlessUUIDField string
}

type AccessSecret = [32]byte
type CertHash = [32]byte

type NodeConnectionInfo struct {
	Endpoint     string
	AccessSecret AccessSecret
	CertHash     CertHash
}

type NodeID int

type NodeConfig struct {
	ClientConfig   ClientConfig
	ConnectionInfo NodeConnectionInfo
}

type Node struct {
	ID            NodeID
	Config        NodeConfig
	CurrentStatus NodeStatus
	TargetStatus  NodeStatus
}

func (s NodeStatus) String() string {
	switch s {
	case NodeStatusStopped:
		return "Stopped"
	case NodeStatusRunning:
		return "Running"
	default:
		return "Unknown"
	}
}

/*type NodeClientCfg struct {
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
}*/
