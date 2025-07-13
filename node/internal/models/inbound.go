package models

type InboundType int

const (
	UnsupportedInbound = iota
	VlessTcpReality
	VlessXHTTP
)

type Inbound struct {
	Tag  string
	Type InboundType
}
