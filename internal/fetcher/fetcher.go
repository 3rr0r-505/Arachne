package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Fetcher struct {
	client *http.Client
}

func NewFetcher(timeout time.Duration) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (f *Fetcher) Fetch(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request for %q: %w", url, err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %q: %w", url, err)
	}

	return resp, nil
}
