package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/html"
)

func main() {
	url := "https://en.wikipedia.org/wiki/Deaths_in_2024#1"

	for {
		fmt.Printf("Crawling page: %s\n", url)
		links, err := crawlPage(url)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
		fmt.Printf("Found %d links. ", len(links))
		url = links[rand.Intn(len(links))]
		fmt.Printf("Visiting:  %s\n", links[0])
	}

	// fmt.Println(links)
}

func crawlPage(pageUrl string) ([]string, error) {
	parsedUrl, err := url.Parse(pageUrl)
	if err != nil {
		return nil, err
	}

	// TODO: put this stuff in its own function
	// First check for robots.txt
	resp, err := http.Get(parsedUrl.Scheme + "://" + parsedUrl.Host + "/robots.txt")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	UNUSED(body)
	// fmt.Printf("Contents of robots.txt:\n%s", string(body))

	resp, err = http.Get(pageUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	links := extractLinks(resp.Body, parsedUrl)
	return links, nil
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
						links = append(links, absoluteLink.String())
					}
				}
			}
		}
	}
}

func UNUSED(x ...interface{}) {}
