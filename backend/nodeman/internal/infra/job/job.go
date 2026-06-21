package job

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/XRay-Addons/xrayman/nodeman/internal/errdefs"
	"go.uber.org/zap"
)

type Op = func(ctx context.Context) error

type Job struct {
	op Op

	interval time.Duration
	cancel   context.CancelFunc
	wg       sync.WaitGroup

	name string
	log  *zap.Logger
}

func NewJob(op Op, jobInterval time.Duration, name string, log *zap.Logger) (*Job, error) {
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
	m := &Job{
		op:       op,
		interval: jobInterval,
		name:     fmt.Sprintf("background job %s", name),
		log:      log,
	}

	return m, nil
}

func (j *Job) Run() error {
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

func (j *Job) Stop() error {
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

func (j *Job) opLoop(ctx context.Context) {
	for {
		if ctx.Err() != nil {
			return
		}

		startTime := time.Now()

		jobCtx, cancel := context.WithTimeout(ctx, j.interval)
		err := j.op(jobCtx)
		cancel()
		if !errors.Is(err, context.Canceled) {
			j.logJobResult(err)
		}

		timeLeft := j.interval - time.Since(startTime)

		select {
		case <-time.After(timeLeft): // immediate if time.left < 0
		case <-ctx.Done():
			return
		}
	}
}

func (j *Job) logJobResult(err error) {
	if err != nil {
		j.log.Error(j.name, zap.Error(err))
		return
	}
}
