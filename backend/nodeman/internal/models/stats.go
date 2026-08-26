package models

type UserStats struct {
	ID       UserID
	Uplink   int64
	Downlink int64
}

type NodePerformance struct {
	OpenConnections int32
	CpuLoad         float32
	RamLoad         float32
	MemLoad         float32
}

type NodeStats struct {
	Users       []UserStats
	Performance NodePerformance
}
