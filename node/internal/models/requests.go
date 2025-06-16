package models

type StartNodeRequest struct {
	Users []User `json:"users" validate:"required"`
}

type AddUsersRequest struct {
	Users []User `json:"users" validate:"required"`
}

type DelUsersRequest struct {
	Users []User `json:"users" validate:"required"`
}
