package models

type User struct {
	Name string `json:"name" validate:"required"`
	UUID string `json:"uuid" validate:"required"`
}
