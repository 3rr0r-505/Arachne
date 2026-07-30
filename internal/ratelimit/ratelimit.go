package ratelimit

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type DomainLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	rps      float64
}

func NewDomainLimiter(rps float64) *DomainLimiter {
	return &DomainLimiter{
		limiters: make(map[string]*rate.Limiter),
		rps:      rps,
	}
}

func (d *DomainLimiter) Wait(ctx context.Context, host string) error {
	if d.rps == 0 {
		return nil
	}

	d.mu.Lock()
	limiter, exist := d.limiters[host]
	if !exist {
		limiter = rate.NewLimiter(rate.Limit(d.rps), 1)
		d.limiters[host] = limiter
	}
	d.mu.Unlock()

	return limiter.Wait(ctx)
}
