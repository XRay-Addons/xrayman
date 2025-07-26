package models

type UserID = int

type User struct {
	ID        UserID
	Name      string
	VlessUUID string
}

type UserStatus int

const (
	UserStatusUnknown UserStatus = iota + 1
	UserDisabled
	UserEnabled
)

func (s UserStatus) String() string {
	switch s {
	case UserDisabled:
		return "Disabled"
	case UserEnabled:
		return "Enabled"
	default:
		return "Unknown"
	}
}
