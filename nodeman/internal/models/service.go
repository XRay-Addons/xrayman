package models

import (
	"github.com/go-faster/jx"
)

type NewNodeParams struct {
	Endpoint  string
	AccessKey AccessKey
}

type NewNodeResult struct {
	Node       Node
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

type DeleteNodeParams struct {
	ID NodeID
}

type DeleteNodeResult struct {
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

type DeleteUserParams struct {
	ID UserID
}

type DeleteUserResult struct {
}

type UserSubParams struct {
	ID   UserID
	Name string
}

type ClientConfigItem = jx.Raw

type UserSubResult struct {
	Expiration    int
	ClientConfigs []ClientConfigItem
}
