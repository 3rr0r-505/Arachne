package scope

import (
	"fmt"
	"net/url"
	"strings"
)

func InScope(seedHost, candidateURL string, allowSubdomains, allowExternal bool) (bool, error) {
	u, err := url.Parse(candidateURL)
	if err != nil {
		return false, fmt.Errorf("invalid candidate URL %q: %w", candidateURL, err)
	}

	if u.Hostname() == seedHost {
		return true, nil
	}

	if strings.HasSuffix(u.Hostname(), "."+seedHost) {
		return allowSubdomains, nil
	}

	return allowExternal, nil
}
