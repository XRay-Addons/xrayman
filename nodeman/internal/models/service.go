package models

type NewNodeParams struct {
	Endpoint  string
	AccessKey AccessKey
}

type NewNodeResult struct {
	ID       NodeID
	Endpoint string
}

type StartNodeParams struct {
	ID NodeID
}

type StartNodeResult struct {
}

type StopNodeParams struct {
	ID NodeID
}

type StopNodeResult struct {
}

type ListNodeParams struct {
}

type ListNodeResult struct {
	Nodes []Node
}

type NewUserParams struct {
	Name string
}

type NewUserResult struct {
	ID          UserID
	Name        string
	UserPageURL string
}

type EnableUserParams struct {
	ID UserID
}

type EnableUserResult struct {
}

type DisableUserParams struct {
	ID UserID
}

type DisableUserResult struct {
}

type ListUserParams struct {
}

type ListUsersResult struct {
	Users []User
}
