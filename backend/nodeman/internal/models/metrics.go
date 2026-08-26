package models

type NodeMetrics struct {
	TotalInbount  int64
	TotalOutbound int64
	CpuLoad       float32
	RamLoad       float32
	MemLoad       float32
}
