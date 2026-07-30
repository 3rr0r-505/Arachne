package crawler

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/3rr0r-505/arachne/internal/config"
	"github.com/3rr0r-505/arachne/internal/fetcher"
	"github.com/3rr0r-505/arachne/internal/ratelimit"
	"github.com/3rr0r-505/arachne/internal/result"
)

func Crawl(ctx context.Context, cfg *config.Config) (<-chan result.PageResult, error) {
	seed, err := url.Parse(cfg.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid seed URL %q: %w", cfg.Url, err)
	}
	seedHost := seed.Hostname()

	jobs := NewJobQueue()
	visited := NewVisitedURLs()
	fetchr := fetcher.NewFetcher(cfg.Timeout, cfg.Retries)
	results := make(chan result.PageResult, cfg.Workers)
	limiter := ratelimit.NewDomainLimiter(cfg.Rate)

	var wg sync.WaitGroup
	var pageCount atomic.Int64

	var cancelCTX context.CancelFunc
	if cfg.CtxTimer > 0 {
		ctx, cancelCTX = context.WithTimeout(ctx, cfg.CtxTimer)
	}

	wg.Add(1)
	jobs.Push(Job{URL: cfg.Url, Depth: 0})

	for id := range cfg.Workers {
		go Worker(id, jobs, results, fetchr, visited, cfg, seedHost, &wg, ctx, limiter, &pageCount)
	}

	go func() {
		wg.Wait()
		jobs.Close()
		close(results)
		if cancelCTX != nil {
			cancelCTX()
		}
	}()

	return results, nil
}
