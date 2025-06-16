package models

type StartNodeRequest struct {
	Users []User `json:"users" validate:"required"`
}

type StartNodeResponse struct {
	NodeProperties Node `json:"node_properties"`
}

type NodeStatusResponse struct {
	NodeStatus NodeStatus `json:"node_status`
}

type AddUsersRequest struct {
	Users []User `json:"users" validate:"required"`
}

type DelUsersRequest struct {
	Users []User `json:"users" validate:"required"`
}
