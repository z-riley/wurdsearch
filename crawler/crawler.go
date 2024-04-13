package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/jimsmart/grobotstxt"
)

const crawlerName = "TurdSeeker"

type Crawler struct {
	frontier *Frontier // queue of links to visit next
	db       Storer    // database abstraction
}

func newCrawler() (*Crawler, error) {

	return &Crawler{
		frontier: nil,
		db:       nil,
	}, nil
}

func (c *Crawler) crawlPage(url *url.URL) error {

	// 1. Get robots.txt
	resp, err := http.Get(url.Scheme + "://" + url.Host + "/robots.txt")
	if err != nil {
		return err
		// TODO: skip robots.txt stuff it one doesn't exist
	}
	defer resp.Body.Close()
	robotsTxt, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// 2. Check we can visit the page
	if !checkPageAllowed(string(robotsTxt), url.String()) {
		return fmt.Errorf("robots.txt forbids visiting page: %s", url)
	}

	// 3. Visit page
	resp, err = http.Get(url.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//  4. Parse page contents
	data, err := parsePage(resp.Body, url)
	if err != nil {
		return err
	}

	// 5. Put the contents in DB
	c.db.StorePlaceholder(data.links)

	return nil
}

func checkPageAllowed(robotsTxt, url string) bool {
	return grobotstxt.AgentAllowed(robotsTxt, crawlerName, url)
}
