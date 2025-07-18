package api

const StartURLPath = "/start"

type StartRequest struct {
	Users []User `json:"users" validate:"required"`
}

const StopURLPath = "/stop"

type StopRequest struct {
}

const StatusURLPath = "/status"

type StatusRequest struct {
}

const EditUsersURLPath = "/users/edit"

type EditUsersRequest struct {
	Add    []User `json:"add" validate:"required"`
	Remove []User `json:"remove" validate:"required"`
}
