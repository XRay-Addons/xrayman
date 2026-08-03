package models

type StartParams struct {
	Users []User
}

type StartResult struct {
	ClientConfigTemplate ClientConfigTemplate
}

type StatusResult struct {
	ServiceStatus ServiceStatus
}

type EditUsersParams struct {
	Add    []User
	Remove []User
}

type StatsResult struct {
	Users []UserStats
}
