package retry

import (
	"testing"

	"github.com/cloudwego/kitex/internal/test"
)

func BenchmarkRandomBackOff_Wait(b *testing.B) {
	min := 10
	max := 100
	bk := &randomBackOff{
		minMS: min,
		maxMS: max,
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ms := int(bk.randomDuration().Milliseconds())
			test.Assert(b, ms <= max && ms >= min, ms)
		}
	})
}

func TestRandomBackOff_Wait(T *testing.T) {
	min := 10
	max := 100
	bk := &randomBackOff{
		minMS: min,
		maxMS: max,
	}
	bk.Wait(1)
}
