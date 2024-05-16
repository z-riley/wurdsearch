package crawler

import (
	"fmt"
	"net/url"
	"testing"
)

func TestCrawlPage(t *testing.T) {
	c, err := NewCrawler(0)
	if err != nil {
		t.Error(err)
	}

	link := "https://urm.wwu.edu/how-create-anchor-jump-link"
	url, err := url.Parse(link)
	if err != nil {
		t.Error(err)
	}
	data, err := c.crawlPage(url)
	if err != nil {
		t.Error(err)
	}

	for _, link := range data.Links {
		fmt.Println(link)
	}
}
