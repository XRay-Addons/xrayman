package service


type PerfCtl interface {
	GetCPUUsage() (int, error)
	GetRAMUsage() (int, error)
}