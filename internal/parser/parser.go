package parser

import (
	"fmt"
	"io"
	"net/url"

	"golang.org/x/net/html"
)

type Form struct {
	Action string
	Method string
}

type ExtractOptions struct {
	IncludeJS    bool
	IncludeForms bool
}

type ExtractResult struct {
	Links   []string
	JSLinks []string
	Forms   []Form
}

func Extract(body io.Reader, baseURL string, opts ExtractOptions) (ExtractResult, error) {
	var result ExtractResult

	base, err := url.Parse(baseURL)
	if err != nil {
		return result, fmt.Errorf("invalid base URL %q: %w", baseURL, err)
	}

	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		if tt == html.ErrorToken {
			if z.Err() != io.EOF {
				return result, fmt.Errorf("html parse error: %w", z.Err())
			}
			break
		}

		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			token := z.Token()
			switch token.Data {
			case "a":
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						ref, err := url.Parse(attr.Val)
						if err != nil {
							continue
						}
						resolved := base.ResolveReference(ref)
						resolved.Fragment = ""
						result.Links = append(result.Links, resolved.String())
					}
				}

			case "script":
				if !opts.IncludeJS {
					continue
				}
				for _, attr := range token.Attr {
					if attr.Key == "src" {
						ref, err := url.Parse(attr.Val)
						if err != nil {
							continue
						}
						resolved := base.ResolveReference(ref)
						resolved.Fragment = ""
						result.JSLinks = append(result.JSLinks, resolved.String())
					}
				}

			case "form":
				if !opts.IncludeForms {
					continue
				}
				var form Form
				for _, attr := range token.Attr {
					switch attr.Key {
					case "action":
						ref, err := url.Parse(attr.Val)
						if err == nil {
							resolved := base.ResolveReference(ref)
							resolved.Fragment = ""
							form.Action = resolved.String()
						}
					case "method":
						form.Method = attr.Val
					}
				}
				result.Forms = append(result.Forms, form)
			}
		}
	}

	return result, nil
}
