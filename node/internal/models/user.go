package models

import "fmt"

type UserID int

type User struct {
	ID        UserID
	Name      string
	VlessUUID string
}

func (u User) VlessEmail() string {
	return fmt.Sprintf("%d-%s", u.ID, u.Name)
}
