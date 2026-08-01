package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"go.uber.org/zap"
)

type Job struct {
	Name    string
	OnStart func(context.Context) error
	OnStop  func(context.Context) error
}

type Closer struct {
	Name   string
	OnStop func(context.Context) error
}

type Lifecycle interface {
	AppendJob(j Job)
	AppendCloser(c Closer)
}

type lifecycle struct {
	log          *zap.Logger
	closeTimeout time.Duration

	jobs    []Job
	closers []Closer
}

func (lc *lifecycle) AppendJob(j Job) {
	lc.jobs = append(lc.jobs, j)
}

func (lc *lifecycle) AppendCloser(c Closer) {
	lc.closers = append(lc.closers, c)
}

func (lc *lifecycle) Run(ctx context.Context) error {
	// wg for waiting all jobs completed (successfully or not)
	wg := sync.WaitGroup{}

	// channel for waiting job error signal (it's signal to stop all)
	errCh := make(chan struct{}, 1)

	// run all jobs
	runErrs := make([]error, len(lc.jobs))
	for idx, job := range lc.jobs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := lc.invokeJobRunner(ctx, job); err != nil {
				runErrs[idx] = err
				select {
				case errCh <- struct{}{}:
				default:
				}
			}
		}()
	}

	select {
	case <-errCh:
	case <-ctx.Done():
	case <-wgChan(&wg):
	}

	// cancel jobs
	closeCtx, cancel := context.WithTimeout(context.Background(), lc.closeTimeout)
	defer cancel()
	closeErr := lc.invokeJobClosers(closeCtx)

	// wait for all jobs completed
	wg.Wait()

	// join all errors
	runErr := xerr.Join(runErrs...)
	return xerr.Join(runErr, closeErr)
}

func (lc *lifecycle) Close() error {
	closeCtx, cancel := context.WithTimeout(context.Background(), lc.closeTimeout)
	defer cancel()
	return lc.invokeClosers(closeCtx)
}

func wgChan(wg *sync.WaitGroup) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	return done
}

func (lc *lifecycle) invokeJobRunner(ctx context.Context, job Job) error {
	if job.OnStart == nil {
		return nil
	}

	lc.log.Info(fmt.Sprintf("run job '%s'...", job.Name))
	err := job.OnStart(ctx)
	if err == nil {
		lc.log.Info(fmt.Sprintf("job '%s' done", job.Name), zap.Error(err))
		return nil
	}

	lc.log.Error(fmt.Sprintf("job '%s'", job.Name), zap.Error(err))
	return xerr.WrapWithInfof(err, "job %s", job.Name)
}

func (lc *lifecycle) invokeJobClosers(ctx context.Context) error {
	errs := make([]error, len(lc.jobs))
	for idx, job := range lc.jobs {
		if job.OnStop == nil {
			continue
		}
		lc.log.Info(fmt.Sprintf("job '%s' stopping signal sent", job.Name))
		if err := job.OnStop(ctx); err != nil {
			lc.log.Error(fmt.Sprintf("job '%s' stoppping", job.Name), zap.Error(err))
			errs[idx] = xerr.WrapWithInfof(err, "job %s", job.Name)
		}
	}
	return xerr.Join(errs...)
}

func (lc *lifecycle) invokeClosers(ctx context.Context) error {
	var errs []error
	for i := len(lc.closers) - 1; i >= 0; i-- {
		closer := lc.closers[i]
		lc.log.Warn(fmt.Sprintf("close '%s'", closer.Name))
		if err := closer.OnStop(ctx); err != nil {
			lc.log.Error(fmt.Sprintf("close '%s'", closer.Name), zap.Error(err))
			errs = append(errs, xerr.WrapWithInfof(err, "closer %s", closer.Name))
		}
	}
	return xerr.Join(errs...)
}
