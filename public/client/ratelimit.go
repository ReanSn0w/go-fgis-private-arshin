package client

import (
	"context"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	interval time.Duration
	next     time.Time
}

func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{interval: interval}
}

func (l *RateLimiter) Wait(ctx context.Context) error {
	if l == nil || l.interval == 0 {
		return nil
	}

	l.mu.Lock()
	now := time.Now()
	waitUntil := l.next
	if now.After(waitUntil) {
		waitUntil = now
	}
	l.next = waitUntil.Add(l.interval)
	l.mu.Unlock()

	timer := time.NewTimer(time.Until(waitUntil))
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
