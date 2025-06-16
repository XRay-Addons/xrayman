package perfctl

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"

	"github.com/XRay-Addons/xrayman/node/internal/errdefs"
)

type PerfCtl struct {
}

// TODO: background goroutine to monitor performance

func (p *PerfCtl) GetCPUUsage() (int, error) {
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

func (p *PerfCtl) GetRAMUsage() (int, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("%w: get cpu usage: %v", errdefs.ErrAccess, err)
	}
	return int(vmStat.UsedPercent + 0.5), nil
}
