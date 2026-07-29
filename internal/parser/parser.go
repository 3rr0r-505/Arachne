package parser

import (
	"fmt"
	"io"
	"net/url"

	"golang.org/x/net/html"
)

func ExtractLinks(body io.Reader, baseURL string) ([]string, error) {
	var urls []string

	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %q: %w", baseURL, err)
	}

	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			if z.Err() != io.EOF {
				return nil, fmt.Errorf("html parse error: %w", z.Err())
			}
			break
		}

		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			token := z.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						ref, err := url.Parse(attr.Val)
						if err != nil {
							continue
						}
						resolved := base.ResolveReference(ref)
						resolved.Fragment = ""
						urls = append(urls, resolved.String())
					}
				}
			}
		}
	}

	return urls, nil
}
