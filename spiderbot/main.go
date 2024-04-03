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
	baseUrl := "http://reddit.com/"
	url, _ := url.Parse(baseUrl)

	fmt.Println(url)

	resp, err := http.Get(baseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	extractLinks(resp.Body)

}

func extractLinks(body io.Reader) []string {
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
						fmt.Println(attr.Val)
						// FIXME: add base url to relative links
						// https://chat.openai.com/c/dc27d85c-002b-47b5-a188-ad68cebae341
					}
				}
			}
		}
	}
}
