package models

type UserName = string

// xray user has Unique Name and (it called 'email' in xray api)
type User struct {
	Name      UserName
	VlessUUID string
}
