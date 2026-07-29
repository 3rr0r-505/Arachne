package crawler

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/3rr0r-505/arachne/internal/config"
	"github.com/3rr0r-505/arachne/internal/fetcher"
	"github.com/3rr0r-505/arachne/internal/result"
)

func Crawl(cfg *config.Config) (<-chan result.PageResult, error) {
	seed, err := url.Parse(cfg.Url)
	if err != nil {
		return nil, fmt.Errorf("invalid seed URL %q: %w", cfg.Url, err)
	}
	seedHost := seed.Hostname()

	jobs := NewJobQueue(cfg.QueueSize)
	visited := NewVisitedURLs()
	fetchr := fetcher.NewFetcher(cfg.Timeout)
	results := make(chan result.PageResult, cfg.Workers)

	var wg sync.WaitGroup

	wg.Add(1)
	jobs.Push(Job{URL: cfg.Url, Depth: 0})

	for id := range cfg.Workers {
		go Worker(id, jobs, results, fetchr, visited, cfg, seedHost, &wg)
	}

	go func() {
		wg.Wait()
		jobs.Close()
		close(results)
	}()

	return results, nil
}
