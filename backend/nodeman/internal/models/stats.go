package models

type UserStats struct {
	ID       UserID
	Upload   int64
	Download int64
}

type NodeStats struct {
	Users []UserStats
}
