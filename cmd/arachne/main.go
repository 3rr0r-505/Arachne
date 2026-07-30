package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/3rr0r-505/arachne/internal/config"
	"github.com/3rr0r-505/arachne/internal/crawler"
	"github.com/3rr0r-505/arachne/internal/validator"
)

const Warning = "\u26A0\uFE0F"
const URL = "\U0001F517\uFE0F"

func main() {
	cfg, err := config.ParseFlags(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := validator.Validate(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancelCTX := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCTX()

	results, err := crawler.Crawl(ctx, cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for r := range results {
		if r.Error != "" {
			fmt.Printf("%s [ERR] %s (depth %d): %s\n", Warning, r.Url, r.Depth, r.Error)
			continue
		}
		fmt.Printf("%s [%d] %s (depth %d)\n", URL, r.Status, r.Url, r.Depth)
		for _, js := range r.JSLinks {
			fmt.Printf("    js: %s\n", js)
		}
		for _, f := range r.Forms {
			fmt.Printf("    form: action=%s method=%s\n", f.Action, f.Method)
		}
	}
}
