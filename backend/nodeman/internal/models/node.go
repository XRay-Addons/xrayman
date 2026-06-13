package models

import (
	"encoding/json"
	"strconv"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/go-faster/jx"
)

type NodeStatus int

const (
	NodeStatusUnknown NodeStatus = iota + 1
	NodeStatusStopped
	NodeStatusRunning
)

type ClientConfigTemplateItem = jx.Raw

type ClientConfigTemplate struct {
	Template        []ClientConfigTemplateItem
	VlessEmailField string
	VlessUUIDField  string
}

func (c ClientConfigTemplate) Value() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", xerr.WrapWithStack(err)
	}
	return string(b), nil
}

func (c *ClientConfigTemplate) Scan(src any) error {
	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return xerr.Newf("unsupported type %T for ClientConfigTemplate", src)
	}
	if err := json.Unmarshal(data, c); err != nil {
		return xerr.WrapWithStack(err)
	}
	return nil
}

type NodeConnectionInfo struct {
	Endpoint  string
	AccessKey AccessKey
}

type NodeID = int

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
