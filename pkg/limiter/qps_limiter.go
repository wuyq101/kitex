package limiter

import (
	"sync/atomic"
	"time"
)

var (
	fixedWindowTime = time.Second
)

// qpsLimiter implements the RateLimiter interface.
type qpsLimiter struct {
	limit    int32
	tokens   int32
	interval time.Duration
	once     int32
	ticker   *time.Ticker
}

// NewQPSLimiter creates qpsLimiter.
func NewQPSLimiter(interval time.Duration, limit int) RateLimiter {
	l := &qpsLimiter{
		limit:    int32(limit),
		tokens:   int32(limit),
		interval: interval,
		once:     calcOnce(interval, limit),
		ticker:   time.NewTicker(interval),
	}
	go l.startTicker()
	return l
}

// UpdateLimit update limitation of QPS. It is **not** concurrent-safe.
func (l *qpsLimiter) UpdateLimit(limit int) {
	once := calcOnce(l.interval, limit)
	atomic.StoreInt32(&l.limit, int32(limit))
	atomic.StoreInt32(&l.once, once)
}

// UpdateQPSLimit update the interval and limit. It is **not** concurrent-safe.
func (l *qpsLimiter) UpdateQPSLimit(interval time.Duration, limit int) {
	atomic.StoreInt32(&l.limit, int32(limit))
	once := calcOnce(interval, limit)
	atomic.StoreInt32(&l.once, once)
	if interval != l.interval {
		l.interval = interval
		l.stopTicker()
		l.ticker = time.NewTicker(interval)
		go l.startTicker()
	}
}

// Acquire adds 1.
func (l *qpsLimiter) Acquire() bool {
	if atomic.LoadInt32(&l.tokens) <= 0 {
		return false
	}
	return atomic.AddInt32(&l.tokens, -1) >= 0
}

// Status returns the current status.
func (l *qpsLimiter) Status() (max, cur int, interval time.Duration) {
	max = int(atomic.LoadInt32(&l.limit))
	cur = int(atomic.LoadInt32(&l.tokens))
	interval = l.interval
	return
}

func (l *qpsLimiter) startTicker() {
	ch := l.ticker.C
	for range ch {
		l.updateToken()
	}
}

func (l *qpsLimiter) stopTicker() {
	l.ticker.Stop()
}

// Some deviation is allowed here to gain better performance.
func (l *qpsLimiter) updateToken() {
	var v int32
	v = atomic.LoadInt32(&l.tokens)
	if v < 0 {
		v = atomic.LoadInt32(&l.once)
	} else if v+atomic.LoadInt32(&l.once) > atomic.LoadInt32(&l.limit) {
		v = atomic.LoadInt32(&l.limit)
	} else {
		v = v + atomic.LoadInt32(&l.once)
	}
	atomic.StoreInt32(&l.tokens, v)
}

func calcOnce(interval time.Duration, limit int) int32 {
	if interval > time.Second {
		interval = time.Second
	}
	once := int32(float64(limit) / (fixedWindowTime.Seconds() / interval.Seconds()))
	if once < 0 {
		once = 0
	}
	return once
}
