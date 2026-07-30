package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Fetcher struct {
	client  *http.Client
	retries int
}

func NewFetcher(timeout time.Duration, retries int) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
		},
		retries: retries,
	}
}

func (f *Fetcher) Fetch(ctx context.Context, url string) (*http.Response, error) {
	var lastErr error

	for attempts := range f.retries + 1 {
		if attempts > 0 {
			delay := time.Duration(1<<attempts) * time.Second

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to build request for %q: %w", url, err)
		}

		resp, err := f.client.Do(req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("failed to fetch %q after %d attempts: %w", url, f.retries+1, lastErr)
}
