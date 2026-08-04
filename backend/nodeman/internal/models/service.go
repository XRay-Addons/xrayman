package models

type NewNodeParams struct {
	Endpoint  string
	AccessKey AccessKey
}

type NewNodeResult struct {
	Node Node
}

type StartNodeParams struct {
	ID NodeID
}

type StopNodeParams struct {
	ID NodeID
}

type ListNodeResult struct {
	Nodes []Node
}

type DeleteNodeParams struct {
	ID NodeID
}

type NewUserParams struct {
	DisplayName string
}

type GetUserParams struct {
	ID   UserID
	Name string
}

type EnableUserParams struct {
	ID UserID
}

type DisableUserParams struct {
	ID UserID
}

type ListUsersResult struct {
	Users []UserView
}

type DeleteUserParams struct {
	ID UserID
}

type UserSubParams struct {
	ID   UserID
	Name string
}

type UserSubResult struct {
	Headers       SubHeaders
	ClientConfigs SubContent
}
