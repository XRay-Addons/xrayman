package models

import (
	"fmt"
	"strconv"
)

type UserID = int

type UserProfile struct {
	ID          UserID
	DisplayName string
	Name        string
	VlessUUID   string
}

func (u UserProfile) VlessEmail() string {
	return fmt.Sprintf("%d-%s", u.ID, u.Name)
}

func (u UserProfile) SubscriptionURL() string {
	return fmt.Sprintf("/sub/%d-%s", u.ID, u.Name)
}

type UserStatus int

const (
	UserStatusUnknown UserStatus = iota + 1
	UserStatusDisabled
	UserStatusEnabled
)

type User struct {
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

func (s UserStatus) StringInt() string {
	return strconv.Itoa(int(s))
}
