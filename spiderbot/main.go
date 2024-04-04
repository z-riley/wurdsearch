package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

func main() {
	baseUrl := "http://levenue.com/"
	baseUrlParsed, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Get(baseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	extractLinks(resp.Body, baseUrlParsed)
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
							log.Fatal(err)
							continue
						}
						// Turn relative links into absolute links
						absoluteLink := baseUrl.ResolveReference(link)
						fmt.Println(absoluteLink)

					}
				}
			}
		}
	}
}
