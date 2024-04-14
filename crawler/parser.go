package main

import (
	"io"
	"net/url"

	"golang.org/x/net/html"
)

type pageData struct {
	links []string
}

func parsePage(body io.Reader, url *url.URL) (pageData, error) {
	return pageData{
		links: extractLinks(body, url),
	}, nil
}

func extractLinks(body io.Reader, baseUrl *url.URL) []string {
	links := []string{}
	z := html.NewTokenizer(body)

	for {
		tokenType := z.Next()
		if tokenType == html.ErrorToken {
			return links // EOF
		} else if tokenType == html.StartTagToken {

			token := z.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link, err := url.Parse(attr.Val)
						if err != nil {
							log.Fatal().Err(err)
							continue
						}
						// Turn relative links into absolute links
						absoluteLink := baseUrl.ResolveReference(link)
						links = append(links, absoluteLink.String())
					}
				}
			}
		}
	}
}
