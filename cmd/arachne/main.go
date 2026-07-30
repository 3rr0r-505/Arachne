package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

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

	results, err := crawler.Crawl(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for res := range results {
		if res.Error != "" {
			// fmt.Printf("[ERR] %s (depth %d): %s\n", res.Url, res.Depth, res.Error)
			fmt.Printf("%s [ERR] %s (depth %d): %s\n", Warning, res.Url, res.Depth, res.Error)
			continue
		}
		// fmt.Printf("[%d] %s (depth %d)\n", res.Status, res.Url, res.Depth)
		fmt.Printf("%s [%d] %s (depth %d)\n", URL, res.Status, res.Url, res.Depth)
	}
}
