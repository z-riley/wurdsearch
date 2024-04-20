package main

import (
	"fmt"
	"net/http"
	"testing"
)

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

	text := extractText(resp.Body)
	fmt.Println(text)
}
