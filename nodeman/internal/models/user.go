package models

type UserID int

type UserProfile struct {
	ID        UserID
	Name      string
	VlessUUID string
}

type UserStatus int

const (
	UserStatusUnknown UserStatus = iota + 1
	UserStatusInactive
	UserStatusActive
)

// Rename it!
type UserTargetState struct {
	User   UserProfile
	Target UserStatus
}

type NodeUsersUpdate struct {
	Add    []UserProfile
	Remove []UserProfile
}

type UserSyncStatus struct {
	User          UserProfile
	TargetStatus  UserStatus
	CurrentStatus UserStatus
}

type UserStatusPatch struct {
	UserID UserID
	Status UserStatus
}

func (s UserStatus) String() string {
	switch s {
	case UserStatusInactive:
		return "Inactive"
	case UserStatusActive:
		return "Active"
	default:
		return "Unknown"
	}
}
