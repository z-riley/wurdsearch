package parser

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"github.com/z-riley/turdsearch/common/store"
	"golang.org/x/net/html"
)

func ParsePage(body io.Reader, url *url.URL, timeAccessed time.Time) (store.PageData, error) {

	// Use the buffer content to create new readers
	bodyContent, err := io.ReadAll(body)
	if err != nil {
		return store.PageData{}, err
	}
	bodyReader1 := bytes.NewReader(bodyContent)
	bodyReader2 := bytes.NewReader(bodyContent)
	bodyReader3 := bytes.NewReader(bodyContent)

	content, err := extractText(bodyReader1)
	if err != nil {
		return store.PageData{}, err
	}

	links := extractLinks(bodyReader2, url)

	return store.PageData{
		Url:          url.String(),
		Title:        extractTitle(bodyReader3),
		LastAccessed: timeAccessed,
		Links:        links,
		Content:      content,
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
						// Only save https links
						if absoluteLink.Scheme != "https" {
							continue
						}
						absoluteLink.Fragment = "" // remove anchor from URL
						links = append(links, absoluteLink.String())
					}
				}
			}
		}
	}
}

func extractText(body io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}

	var textContent strings.Builder
	doc.Find("p, h1, h2, h3, h4, h5, h6, li").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		// Replace multiple newlines and trim space
		re := regexp.MustCompile(`\s+`)
		cleanedText := re.ReplaceAllString(text, " ")
		cleanedText = strings.TrimSpace(cleanedText)

		if cleanedText != "" {
			textContent.WriteString(cleanedText + "\n")
		}
	})
	return textContent.String(), nil
}

// extractTitle extracts the title from the header of an HTML page using goquery.
func extractTitle(body io.Reader) string {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "no title found"
	}
	title := doc.Find("title").Text()
	if title == "" {
		return "title not found"
	}
	fmt.Println(title)
	return strings.TrimSpace(title)
}
