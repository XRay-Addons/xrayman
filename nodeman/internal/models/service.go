package models

import "encoding/json"

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
	VisibleName string
}

type NewUserResult struct {
	ID          UserID
	VisibleName string
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

type GetUserSubParams struct {
	ID   UserID
	Name string
}

type Subscription = json.RawMessage

type GetUserSubResult = []Subscription
