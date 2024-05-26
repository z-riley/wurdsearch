package parser

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestParsePage(t *testing.T) {
	URL := "https://example.com/"
	resp, err := http.Get(URL)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	parsedURL, err := url.Parse(URL)
	if err != nil {
		t.Fatal(err)
	}

	pageData, err := ParsePage(resp.Body, parsedURL, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", pageData)

}

func TestExtractText(t *testing.T) {

	// URLs for manual testing
	url := "https://www.vinted.com/"
	url = "https://www.reddit.com/r/linuxquestions/comments/jb981q/zsh_vs_fish/"
	url = "https://www.youtube.com/"
	url = "https://en.wikipedia.org//wiki/Main_Page"
	url = "https://en.wikipedia.org/wiki/V%C3%A4dersolstavlan"
	url = "https://alwaysjudgeabookbyitscover.com/"
	url = "https://example.com/"

	resp, err := http.Get(url)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	text, err := extractText(resp.Body)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(text)
}
