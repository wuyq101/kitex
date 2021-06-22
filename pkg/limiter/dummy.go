package limiter

import "time"

// DummyConcurrencyLimiter implements ConcurrencyLimiter but without actual limitation.
type DummyConcurrencyLimiter struct{}

// Acquire .
func (dcl *DummyConcurrencyLimiter) Acquire() bool { return true }

// Release .
func (dcl *DummyConcurrencyLimiter) Release() {}

// Status .
func (dcl *DummyConcurrencyLimiter) Status() (limit, occupied int) { return }

// DummyRateLimiter implements RateLimiter but without actual limitation.
type DummyRateLimiter struct{}

// Acquire .
func (drl *DummyRateLimiter) Acquire() bool { return true }

// Status .
func (drl *DummyRateLimiter) Status() (max, current int, interval time.Duration) { return }
