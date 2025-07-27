package models

type NodeStatus int

const (
	NodeStatusUnknown NodeStatus = iota + 1
	NodeStatusStopped
	NodeStatusRunning
)

type ClientTemplate struct {
	Config         string
	UsernameField  string
	VlessUUIDField string
}

type NodeConfig struct {
	ClientTemplate ClientTemplate
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

/*type NodeID = int

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
}*/
