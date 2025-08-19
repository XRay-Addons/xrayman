package models

type UserID int

type UserProfile struct {
	VisibleName string
	Name        string
	VlessUUID   string
}

type UserStatus int

const (
	UserStatusUnknown UserStatus = iota + 1
	UserStatusDisabled
	UserStatusEnabled
)

type User struct {
	ID           UserID
	Profile      UserProfile
	TargetStatus UserStatus
}

type UserSyncStatus struct {
	User          User
	CurrentStatus UserStatus
}

type UserStatusPatch struct {
	UserID UserID
	Status UserStatus
}

type NodeUsersUpdate struct {
	Add    []UserProfile
	Remove []UserProfile
}

func (s UserStatus) String() string {
	switch s {
	case UserStatusEnabled:
		return "Enabled"
	case UserStatusDisabled:
		return "Disabled"
	default:
		return "Unknown"
	}
}

/*type UserProfile struct {
	ID        UserID
	Name      string
	VlessUUID string
}

type UserStatus int

const (
	UserStatusUnknown UserStatus = iota + 1
	UserStatusDisabled
	UserStatusEnabled
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
	case UserStatusDisabled:
		return "Inactive"
	case UserStatusEnabled:
		return "Active"
	default:
		return "Unknown"
	}
}*/
