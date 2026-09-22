package xerrgroup

import (
	"github.com/XRay-Addons/xrayman/common/safego"
	"github.com/XRay-Addons/xrayman/common/xerr"
	"golang.org/x/sync/errgroup"
)

type Group struct {
	g     errgroup.Group
	tasks []func() error
}

func (g *Group) Go(fn func() error) {
	g.tasks = append(g.tasks, fn)
}

// return joint error and task errors in order tasks were added
func (g *Group) Wait() (error, []error) {
	errs := make([]error, len(g.tasks))
	for idx, task := range g.tasks {
		g.g.Go(func() error {
			errs[idx] = safego.Invoke(task)
			return nil
		})
	}
	/*err*/ _ = g.g.Wait() // err is always nil, all errors in errs
	return xerr.Join(errs...), errs
}
