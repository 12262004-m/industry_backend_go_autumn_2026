package main

import (
	"sync"
	"time"
)

type Clock interface{ Now() time.Time }
type Limiter struct {
	mu     sync.Mutex
	clock  Clock
	rate   float64
	burst  int
	tokens float64
	last   time.Time
}

func NewLimiter(clock Clock, ratePerSec float64, burst int) *Limiter {
	l := &Limiter{
		clock: clock,
		rate:  ratePerSec,
		burst: burst,
	}

	if burst > 0 && clock != nil {
		l.tokens = float64(burst)
		l.last = clock.Now()
	}

	return l
}

func (l *Limiter) AllowN(n int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if n <= 0 || n > l.burst || l.burst <= 0 || l.clock == nil {
		return false
	}

	now := l.clock.Now()
	if l.last.Before(now) {
		prev := now.Sub(l.last).Seconds()
		l.tokens = min(float64(l.burst), l.tokens+prev*max(l.rate, 0))
		l.last = now
	}

	if l.tokens < float64(n) {
		return false
	}

	l.tokens -= float64(n)
	return true
}
