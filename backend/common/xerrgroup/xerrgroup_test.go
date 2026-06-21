package xerrgroup

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/XRay-Addons/xrayman/common/xerr"
)

func TestGroup_DemonstrateRaceCondition(t *testing.T) {
	g, ctx := WithContext(context.Background())
	_ = ctx

	for i := 0; i < 10; i++ {

		taskNum := i

		g.Go(func() error {
			time.Sleep(10 * time.Millisecond)
			return xerr.Newf("error task #%d", taskNum)
		})
	}

	err := g.Wait()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func TestGroup_GoroutinePanic(t *testing.T) {
	done := make(chan struct{})

	go func() {
		defer close(done)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 5; i++ {
				fmt.Println("[worker-1] tick", i)
				time.Sleep(300 * time.Millisecond)
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < 3; i++ {
				fmt.Println("[worker-2] tick", i)
				time.Sleep(200 * time.Millisecond)
			}

			panic("boom in worker-2")
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 5; i++ {
				fmt.Println("[worker-3] tick", i)
				time.Sleep(250 * time.Millisecond)
			}
		}()

		wg.Wait()
		fmt.Println("goroutines finished normally")
	}()

	<-done
}
