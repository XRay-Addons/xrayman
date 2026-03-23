package models

type StartParams struct {
	Users []User
}

type StartResult struct {
	ClientConfigTemplate ClientConfigTemplate
}

type StopParams struct {
}

type StopResult struct {
}

type StatusParams struct {
}

type StatusResult struct {
	ServiceStatus ServiceStatus
}

type EditUsersParams struct {
	Add    []User
	Remove []User
}

type EditUsersResult struct {
}
