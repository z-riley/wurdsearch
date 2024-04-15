package main

import (
	"io"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

type pageData struct {
	url          string    `bson:"url"`
	lastAccessed time.Time `bson:"lastAccessed"`
	links        []string  `bson:"links"`
	content      string    `bson:"content"`
}

func parsePage(body io.Reader, url *url.URL, timeAccessed time.Time) (pageData, error) {
	return pageData{
		url:          url.String(),
		lastAccessed: timeAccessed,
		links:        extractLinks(body, url),
		content:      extractContent(body),
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

func extractContent(body io.Reader) string {
	return "example"
}
