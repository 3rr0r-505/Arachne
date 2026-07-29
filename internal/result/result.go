package result

import "time"

type PageResult struct {
	Url       string
	Status    int
	Depth     int
	TimeStamp time.Time
	Error     string
}

func NewPage(url string, depth int) PageResult {
	return PageResult{
		Url:       url,
		Depth:     depth,
		TimeStamp: time.Now(),
	}
}
