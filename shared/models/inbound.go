package models

type InboundType int

const (
	UnsupportedInbound InboundType = iota
	VlessTcpReality                = iota
	VlessXHTTP
)

type Inbound struct {
	Tag  string
	Type InboundType
}
