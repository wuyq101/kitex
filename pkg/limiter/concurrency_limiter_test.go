package limiter

import (
	"sync"
	"testing"
	"time"

	"github.com/cloudwego/kitex/internal/test"
)

func TestConLimiter(t *testing.T) {
	lim := NewConcurrencyLimiter(50)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			_, now := lim.Status()
			test.Assert(t, now <= 50)
			time.Sleep(time.Millisecond / 10)
		}
	}()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				if lim.Acquire() {
					time.Sleep(time.Millisecond)
					lim.Release()
				} else {
					time.Sleep(time.Millisecond)
				}
			}
		}()
	}

	wg.Wait()
}
