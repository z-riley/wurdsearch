package parser

import (
	"bytes"
	"io"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/store"
	"golang.org/x/net/html"
)

func ParsePage(body io.Reader, url *url.URL, timeAccessed time.Time) store.PageData {
	var buf bytes.Buffer
	tee := io.TeeReader(body, &buf)

	content := ensureUTF8(extractText(tee), ';')
	links := extractLinks(&buf, url)

	return store.PageData{
		Url:          url.String(),
		LastAccessed: timeAccessed,
		Links:        links,
		Content:      content,
	}
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

func extractText(body io.Reader) string {
	var text string
	z := html.NewTokenizer(body)

	token := z.Token()

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return text // EOF
		case tt == html.StartTagToken:
			token = z.Token()
		case tt == html.TextToken:
			if token.Data == "script" || token.Data == "style" {
				continue
			}
			content := strings.TrimSpace(html.UnescapeString(string(z.Text())))
			if len(content) > 0 {
				text = text + content + " "
			}
		}
	}
}

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
