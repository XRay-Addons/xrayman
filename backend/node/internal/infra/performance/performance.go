package performance

import (
	"context"
	"os"
	"sync"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/XRay-Addons/xrayman/node/internal/models"
	"github.com/XRay-Addons/xrayman/node/internal/service"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

type Performance struct {
	mu   sync.Mutex
	proc *process.Process
}

var _ service.Performance = (*Performance)(nil)

func New() (*Performance, error) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return nil, xerr.WrapWithInfo(err, "get self addr")
	}
	return &Performance{proc: proc}, nil
}

func (p *Performance) GetPerformance(ctx context.Context) (*models.NodePerformance, error) {
	if !p.mu.TryLock() {
		return nil, xerr.New("performance fetcher locked")
	}
	defer p.mu.Unlock()

	var result models.NodePerformance

	cpuPercents, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	if len(cpuPercents) > 0 {
		result.CpuLoad = float32(cpuPercents[0])
	}

	vm, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	result.RamLoad = float32(vm.UsedPercent)

	diskStat, err := disk.UsageWithContext(ctx, "/")
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	result.MemLoad = float32(diskStat.UsedPercent)

	conns, err := net.ConnectionsPidWithContext(ctx, "all", p.proc.Pid)
	if err != nil {
		return nil, xerr.WrapWithStack(err)
	}
	result.OpenConnections = int32(len(conns))

	return &result, nil
}
