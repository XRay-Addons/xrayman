package xerrgroup

import (
	"context"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"golang.org/x/sync/errgroup"
)

type Group struct {
	g     *errgroup.Group
	tasks []func() error
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	return &Group{g: g}, ctx
}

func (g *Group) Go(fn func() error) {
	g.tasks = append(g.tasks, fn)
}

func (g *Group) Wait() error {
	errs := make([]error, len(g.tasks))
	for idx, task := range g.tasks {
		g.g.Go(func() error {
			defer func() {
				if r := recover(); r != nil {
					errs[idx] = xerr.Newf("panic: %v", r)
				}
			}()

			errs[idx] = task()
			return nil
		})
	}
	err := g.g.Wait()
	errs = append(errs, xerr.WrapWithStack(err))
	return xerr.Join(errs...)
}
