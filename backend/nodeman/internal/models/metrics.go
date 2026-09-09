package models

type NodeTrafficMetrics struct {
	Upload   int64
	Download int64
}

type NodeMetrics struct {
	ID          NodeID
	Endpoint    string
	Traffic     NodeTrafficMetrics
	Performance NodePerformance
}
