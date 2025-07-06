package models

type NodeStatus string

const (
	NodeRunning NodeStatus = "Running"
	NodeStopped NodeStatus = "Stopped"
)
