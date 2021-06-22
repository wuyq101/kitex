// +build !race

package limiter

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cloudwego/kitex/internal/test"
)

func TestQpsLimiter(t *testing.T) {
	limiter := NewQPSLimiter(time.Second/10, 1000)
	for i := 0; i < 1000; i++ {
		test.Assert(t, limiter.Acquire())
	}
	for i := 0; i < 1000; i++ {
		test.Assert(t, !limiter.Acquire())
	}

	time.Sleep(time.Second/2 + 10*time.Millisecond)

	for i := 0; i < 500; i++ {
		test.Assert(t, limiter.Acquire())
	}
	for i := 0; i < 500; i++ {
		test.Assert(t, !limiter.Acquire())
	}

	var wg sync.WaitGroup
	var success int32
	limiter.(*qpsLimiter).UpdateQPSLimit(time.Second/5, 10000)
	max, _, interval := limiter.Status()
	test.Assert(t, max == 10000)
	test.Assert(t, interval == time.Second/5)
	time.Sleep(time.Second)
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok := limiter.Acquire()
			if ok {
				atomic.AddInt32(&success, 1)
			}
		}()
	}
	wg.Wait()
	if success > 11000 {
		t.Log(success)
		t.Fail()
	}
	t.Log(success)
	// test negative limit
	limiter = NewQPSLimiter(time.Second, -100)
	time.Sleep(time.Second * 2)
	tokens := atomic.LoadInt32(&limiter.(*qpsLimiter).tokens)
	once := atomic.LoadInt32(&limiter.(*qpsLimiter).once)
	test.Assert(t, tokens == once)
}

func Test_calOnce(t *testing.T) {
	ret := calcOnce(10*time.Millisecond, 100)
	test.Assert(t, ret > 0)

	// interval > time.Second
	ret = calcOnce(time.Minute, 100)
	test.Assert(t, ret == 100)

	// limit < 0
	ret = calcOnce(time.Second, -100)
	test.Assert(t, ret == 0)
}
