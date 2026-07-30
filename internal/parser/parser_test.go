package parser

import (
	"strings"
	"testing"
)

func TestExtractor(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		baseURL string
		want    []string
	}{
		{
			name:    "absolute href",
			html:    `<a href="https://example.com/page">link</a>`,
			baseURL: "https://example.com",
			want:    []string{"https://example.com/page"},
		},
		{
			name:    "relative href resolved against base",
			html:    `<a href="/about">link</a>`,
			baseURL: "https://example.com/blog/post1",
			want:    []string{"https://example.com/about"},
		},
		{
			name:    "fragment stripped",
			html:    `<a href="https://example.com/page#section">link</a>`,
			baseURL: "https://example.com",
			want:    []string{"https://example.com/page"},
		},
		{
			name:    "multiple anchor tags",
			html:    `<a href="/a">one</a><a href="/b">two</a>`,
			baseURL: "https://example.com",
			want:    []string{"https://example.com/a", "https://example.com/b"},
		},
		{
			name:    "anchor with no href",
			html:    `<a class="nolink">no href here</a>`,
			baseURL: "https://example.com",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Extract(strings.NewReader(tt.html), tt.baseURL, ExtractOptions{})
			if err != nil {
				t.Fatalf("Extract() error = %v", err)
			}

			if len(got.Links) != len(tt.want) {
				t.Fatalf("Extract().Links = %v, want %v", got.Links, tt.want)
			}
			for i := range got.Links {
				if got.Links[i] != tt.want[i] {
					t.Errorf("Extract().Links[%d] = %v, want %v", i, got.Links[i], tt.want[i])
				}
			}
		})
	}
}
