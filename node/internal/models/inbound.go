package models

type InboundType int

const (
	VlessTcpReality = iota
	VlessXHTTP
)

type Inbound struct {
	Tag  string
	Type InboundType
}
