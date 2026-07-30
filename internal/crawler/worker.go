package crawler

import (
	"context"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/3rr0r-505/arachne/internal/config"
	"github.com/3rr0r-505/arachne/internal/fetcher"
	"github.com/3rr0r-505/arachne/internal/parser"
	"github.com/3rr0r-505/arachne/internal/ratelimit"
	"github.com/3rr0r-505/arachne/internal/result"
	"github.com/3rr0r-505/arachne/internal/scope"
)

func Worker(
	id int,
	jobs *JobQ,
	results chan<- result.PageResult,
	fetchr *fetcher.Fetcher,
	visited *VisitedURLs,
	cfg *config.Config,
	seedHost string,
	wg *sync.WaitGroup,
	ctx context.Context,
	limiter *ratelimit.DomainLimiter,
	pageCount *atomic.Int64,
) {
	for {
		job, ok := jobs.Pop()
		if !ok {
			return
		}
		func() {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
			}

			if cfg.Depth != -1 && job.Depth > cfg.Depth {
				return
			}

			if cfg.MaxPages > 0 && pageCount.Load() >= int64(cfg.MaxPages) {
				return
			}
			defer pageCount.Add(1)

			u, err := url.Parse(job.URL)
			if err == nil {
				if err := limiter.Wait(ctx, u.Hostname()); err != nil {
					results <- result.PageResult{
						Url:       job.URL,
						Depth:     job.Depth,
						TimeStamp: time.Now(),
						Error:     err.Error(),
					}
					return
				}
			}

			resp, err := fetchr.Fetch(ctx, job.URL)
			if err != nil {
				results <- result.PageResult{
					Url:       job.URL,
					Depth:     job.Depth,
					TimeStamp: time.Now(),
					Error:     err.Error(),
				}
				return
			}

			extracted, err := parser.Extract(resp.Body, job.URL, parser.ExtractOptions{
				IncludeJS:    cfg.JS,
				IncludeForms: cfg.Forms,
			})
			resp.Body.Close()
			if err != nil {
				results <- result.PageResult{
					Url:       job.URL,
					Status:    resp.StatusCode,
					Depth:     job.Depth,
					TimeStamp: time.Now(),
					Error:     err.Error(),
				}
				return
			}

			for _, link := range extracted.Links {
				inscope, err := scope.InScope(seedHost, link, cfg.SubDomains, cfg.External)
				if err != nil || !inscope {
					continue
				}
				if visited.MarkIfNew(link) {
					wg.Add(1)
					jobs.Push(Job{URL: link, Depth: job.Depth + 1})
				}
			}

			results <- result.PageResult{
				Url:       job.URL,
				Status:    resp.StatusCode,
				Depth:     job.Depth,
				TimeStamp: time.Now(),
				Error:     "",
				JSLinks:   extracted.JSLinks,
				Forms:     extracted.Forms,
			}
		}()
	}
}
