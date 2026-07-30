package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Fetcher struct {
	client  *http.Client
	retries int
}

func NewFetcher(timeout time.Duration, retries int, proxy string) (*Fetcher, error) {
	client := &http.Client{Timeout: timeout}

	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL %q: %w", proxy, err)
		}

		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	return &Fetcher{client: client, retries: retries}, nil
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
