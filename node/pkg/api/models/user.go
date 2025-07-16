package models

type User struct {
	ID        int    `json:"id" validate:"required"`
	Name      string `json:"name" validate:"required"`
	VlessUUID string `json:"vlessUuid" validate:"required"`
}
