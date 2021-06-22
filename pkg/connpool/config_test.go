package connpool

import (
	"testing"
	"time"

	"github.com/cloudwego/kitex/internal/test"
)

func TestCheckPoolConfig(t *testing.T) {
	cfg := IdleConfig{MaxIdleTimeout: 0}
	test.Assert(t, CheckPoolConfig(cfg).MaxIdleTimeout == defaultMaxIdleTimeout)
	cfg = IdleConfig{MaxIdleTimeout: time.Millisecond}
	test.Assert(t, CheckPoolConfig(cfg).MaxIdleTimeout == minMaxIdleTimeout)
}
