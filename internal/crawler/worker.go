package crawler

import (
	"context"
	"sync"
	"time"

	"github.com/3rr0r-505/arachne/internal/config"
	"github.com/3rr0r-505/arachne/internal/fetcher"
	"github.com/3rr0r-505/arachne/internal/parser"
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

			links, err := parser.ExtractLinks(resp.Body, job.URL)
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

			for _, link := range links {
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
			}
		}()
	}
}
