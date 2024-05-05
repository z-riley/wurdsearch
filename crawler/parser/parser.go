package parser

import (
	"bytes"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/common/store"
	"golang.org/x/net/html"
)

func ParsePage(body io.Reader, url *url.URL, timeAccessed time.Time) (store.PageData, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(body, &buf)

	content, err := extractText(tee)
	if err != nil {
		return store.PageData{}, nil
	}

	links := extractLinks(&buf, url)

	return store.PageData{
		Url:          url.String(),
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

// UNUSED:
// ensureUTF8 replaces multi-byte characters with a replacement character
func ensureUTF8(input string, replacement rune) string {
	if utf8.ValidString(input) {
		return input
	}
	v := make([]rune, 0, len(input))
	for i, r := range input {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(input[i:])
			if size == 1 {
				v = append(v, replacement) // replace invalid rune
				continue
			}
		}
		v = append(v, r)
	}
	return string(v)
}
