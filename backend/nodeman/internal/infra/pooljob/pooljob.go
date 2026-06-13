package pooljob

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"github.com/XRay-Addons/xrayman/nodeman/internal/models"
	"go.uber.org/zap"
)

type PoolJob struct {
	op PoolOp

	interval time.Duration
	cancel   context.CancelFunc
	wg       sync.WaitGroup

	name string
	log  *zap.Logger
}

func New(op PoolOp, jobInterval time.Duration, name string, log *zap.Logger) (*PoolJob, error) {
	if op == nil {
		return nil, errdefs.NilArg("op")
	}
	if jobInterval == 0 {
		return nil, errdefs.NilArg("jobInterval")
	}
	if log == nil {
		return nil, errdefs.NilArg("log")
	}
	// init default options
	m := &PoolJob{
		op:       op,
		interval: jobInterval,
		name:     fmt.Sprintf("background pool job %s", name),
		log:      log,
	}

	return m, nil
}

func (j *PoolJob) Run() error {
	if j == nil {
		return errdefs.NilCall()
	}

	ctx, cancel := context.WithCancel(context.Background())
	j.cancel = cancel

	// run op loop
	j.wg.Add(1)
	defer j.wg.Done()
	j.opLoop(ctx)

	return ctx.Err() //nolint:wrapcheck
}

func (j *PoolJob) Stop() error {
	if j == nil {
		return nil
	}
	if j.cancel != nil {
		j.cancel()
		j.cancel = nil
	}
	j.wg.Wait()
	return nil
}

func (j *PoolJob) opLoop(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		// set job time limit to m.interval
		jobCtx, cancel := context.WithTimeout(ctx, j.interval)
		jobRes, err := j.op(jobCtx)
		cancel()
		j.logJobResult(jobRes, err)

		select {
		case <-time.After(j.interval):
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (j *PoolJob) logJobResult(r *models.PoolOpResult, err error) {
	if err != nil {
		j.log.Error(j.name, zap.Error(err))
		return
	}
	for _, n := range r.Nodes {
		if n.Err == nil {
			j.log.Info(fmt.Sprintf("%s OK", j.name),
				zap.String("nodeID", strconv.Itoa(n.ID)),
				zap.String("endpoint", n.Endpoint))
		} else {
			j.log.Error(j.name,
				zap.Error(n.Err),
				zap.String("nodeID", strconv.Itoa(n.ID)),
				zap.String("endpoint", n.Endpoint))
		}
	}
}
