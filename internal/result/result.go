package result

import (
	"time"

	"github.com/3rr0r-505/arachne/internal/parser"
)

type PageResult struct {
	Url       string
	Status    int
	Depth     int
	TimeStamp time.Time
	Error     string
	JSLinks   []string
	Forms     []parser.Form
}

func NewPage(url string, depth int) PageResult {
	return PageResult{
		Url:       url,
		Depth:     depth,
		TimeStamp: time.Now(),
	}
}
