package models

type NodeTraffic struct {
	Upload   int64
	Download int64
}

type NodeMetrics struct {
	ID       NodeID
	Endpoint string
	Traffic  NodeTraffic
	Perf     NodePerformance
}
