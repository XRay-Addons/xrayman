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

type NodeConnectionInfo struct {
	Endpoint  string
	AccessKey string
}

type NodeID int

type NodeConfig struct {
	ID             NodeID
	ClientConfig   ClientConfig
	ConnectionInfo NodeConnectionInfo
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
