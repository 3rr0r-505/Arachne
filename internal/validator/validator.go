package validator

import (
	"fmt"
	"strings"

	"github.com/3rr0r-505/arachne/internal/config"
)

func Validate(cfg *config.Config) error {
	var invalid []string

	if cfg.Url == "" {
		invalid = append(invalid, "url is required")
	}

	if cfg.Depth < 0 {
		invalid = append(invalid, "depth must be >= 0")
	}

	if cfg.Workers <= 0 {
		invalid = append(invalid, "workers must be > 0")
	}

	if cfg.MaxPages < 0 {
		invalid = append(invalid, "max-pages must be >= 0")
	}

	if cfg.QueueSize <= 0 {
		invalid = append(invalid, "queue-size must be > 0")
	}

	if cfg.Timeout <= 0 {
		invalid = append(invalid, "timeout must be > 0")
	}

	if len(invalid) > 0 {
		return fmt.Errorf("invalid config: \n%s", strings.Join(invalid, "\n"))
	}

	return nil
}
