package models

type UserStats struct {
	ID       UserID
	Uplink   int64
	Downlink int64
}
