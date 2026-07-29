package crawler

import "sync"

type VisitedURLs struct {
	mu      sync.Mutex
	visited map[string]bool
}

func NewVisitedURLs() *VisitedURLs {
	return &VisitedURLs{
		visited: make(map[string]bool),
	}
}

func (v *VisitedURLs) MarkIfNew(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.visited[url] {
		return false
	}

	v.visited[url] = true
	return true
}
