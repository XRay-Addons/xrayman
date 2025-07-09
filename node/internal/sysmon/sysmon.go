package perfctl

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
	"github.com/XRay-Addons/xrayman/node/internal/models"
)

type SysMon struct {
}

// TODO: background goroutine to monitor performance

func (p *SysMon) GetSystemUsage() (*models.SysStat, error) {
	cpu, err := getCPUUsage()
	if err != nil {
		return nil, err
	}
	ram, err := getRAMUsage()
	if err != nil {
		return nil, err
	}

	return &models.SysStat{
		CPULoad: cpu,
		RAMLoad: ram,
	}, nil
}

func getCPUUsage() (int, error) {
	measureInterval := 100 * time.Millisecond
	percentages, err := cpu.Percent(measureInterval, false)
	if err != nil {
		return 0, fmt.Errorf("%w: get cpu usage: %v", errdefs.ErrAccess, err)
	}
	if len(percentages) == 0 {
		return 0, fmt.Errorf("%w: no CPU data", errdefs.ErrAccess)
	}
	return int(percentages[0] + 0.5), nil
}

func getRAMUsage() (int, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("%w: get cpu usage: %v", errdefs.ErrAccess, err)
	}
	return int(vmStat.UsedPercent + 0.5), nil
}
