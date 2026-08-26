package models

type StartParams struct {
	Users []User
}

type StartResult struct {
	ClientConfigTemplate ClientConfigTemplate
	Version              string
}

type StatusResult struct {
	ServiceStatus ServiceStatus
}

type EditUsersParams struct {
	Add    []User
	Remove []User
}

type NodePerformance struct {
	OpenConnections int32
	CpuLoad         float32
	RamLoad         float32
	MemLoad         float32
}

type StatsResult struct {
	Users       []UserStats
	Performance NodePerformance
}
