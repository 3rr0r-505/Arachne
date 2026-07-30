package config

import (
	"errors"
	"flag"
	"fmt"
	"time"
)

type Config struct {
	Url                   string
	Depth                 int
	Workers               int
	Timeout               time.Duration
	MaxPages              int
	SubDomains            bool
	External              bool
	Rate                  float64
	Retries               int
	CtxTimer              time.Duration
	Robots, Forms, JS     bool
	Format, Output, Proxy string
}

func ParseFlags(args []string) (*Config, error) {
	fs := flag.NewFlagSet("arachne", flag.ContinueOnError)

	url := fs.String("url", "", "seed URL to start crawling (required)")
	depth := fs.Int("depth", -1, "max crawl depth (-1 = unlimited)")
	workers := fs.Int("workers", 10, "number of concurrent worker goroutines")
	timeout := fs.Duration("timeout", 10*time.Second, "per-request timeout")

	maxPages := fs.Int("max-pages", 0, "hard cap on total pages crawled (0 = unlimited)")
	subDomains := fs.Bool("subs", false, "include subdomains as in-scope")
	external := fs.Bool("external", false, "follow external links too")

	rate := fs.Float64("rate", 0, "requests/sec per domain")
	retries := fs.Int("retries", 2, "max retry attempts on failed fetch")
	ctxTimer := fs.Duration("ctx-timeout", 0, "overall crawl timeout")

	robots := fs.Bool("robots", true, "respect robots.txt")
	forms := fs.Bool("forms", false, "extract forms found on pages")
	js := fs.Bool("js", false, "extract JS file links")

	format := fs.String("format", "text", "output format: text | json | csv")
	output := fs.String("out", "", "output file path (default stdout)")
	proxy := fs.String("proxy", "", "proxy URL to route requests through")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil, err // pass through untouched, let caller decide (exit 0)
		}
		return nil, fmt.Errorf("invalid flags: %w", err)
	}

	return &Config{
		Url:        *url,
		Depth:      *depth,
		Workers:    *workers,
		Timeout:    *timeout,
		MaxPages:   *maxPages,
		SubDomains: *subDomains,
		External:   *external,
		Rate:       *rate,
		Retries:    *retries,
		CtxTimer:   *ctxTimer,
		Robots:     *robots,
		Forms:      *forms,
		JS:         *js,
		Format:     *format,
		Output:     *output,
		Proxy:      *proxy,
	}, nil
}
