package models

type UserStats struct {
	ID       UserID
	Uplink   int64
	Downlink int64
}

type NodeStats struct {
	Users []UserStats
}
