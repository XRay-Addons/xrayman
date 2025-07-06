package models

type User struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	UUID string `json:"uuid" validate:"required"`
}
