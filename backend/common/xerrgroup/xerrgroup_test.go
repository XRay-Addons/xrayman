package xerrgroup

import (
	"fmt"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
	"github.com/stretchr/testify/require"
)

func TestGroup(t *testing.T) {
	g := Group{}
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			time.Sleep(10 * time.Millisecond)
			return xerr.Newf("error task #%d", i)
		})
	}

	err, errs := g.Wait()
	require.Error(t, err)
	require.Equal(t, 10, len(errs))
	for i := range 10 {
		require.Error(t, errs[i])
	}
}
func TestGroup_Panic(t *testing.T) {
	g := Group{}
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			time.Sleep(10 * time.Millisecond)
			panic(fmt.Sprintf("error task #%d", i))
		})
	}

	err, errs := g.Wait()
	require.Error(t, err)
	require.Equal(t, 10, len(errs))
	for i := range 10 {
		require.Error(t, errs[i])
	}
}
