package models

import (
	"strconv"
)

type NodeStatus int

const (
	NodeStatusUnknown NodeStatus = iota + 1
	NodeStatusStopped
	NodeStatusRunning
)

type ClientConfigTemplateItem = string

type ClientConfigTemplate struct {
	Template        []ClientConfigTemplateItem
	VlessEmailField string
	VlessUUIDField  string
}

type NodeConnectionInfo struct {
	Endpoint  string
	AccessKey AccessKey
}

type NodeID int

type NodeConfig struct {
	ClientConfigTemplate ClientConfigTemplate
	ConnectionInfo       NodeConnectionInfo
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

func (s NodeStatus) StringInt() string {
	return strconv.Itoa(int(s))
}

/*type NodeClientConfig struct {
	Template       string
	UserNameField  string
	VlessUUIDField string
}

type NodeProperties struct {
	ClientConfig NodeClientConfig
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
