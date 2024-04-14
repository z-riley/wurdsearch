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
		frontier: newFrontier(),
		db:       newStore(),
	}, nil
}

func (c *Crawler) crawlForever() error {

	for {
		// 1. Get next link from frontier
		link, err := c.frontier.queue.Dequeue()
		if err != nil {
			return err
		}
		link, ok := link.(string)
		if !ok {
			return fmt.Errorf("%v is not castable to string", link)
		}

		// 2. Crawl page
		url, err := url.Parse(link.(string))
		if err != nil {
			return err
		}
		data, err := c.crawlPage(url)
		if err != nil {
			return err
		}
		fmt.Print("Data:")
		fmt.Println(data)

		// 3. Put new links into frontier

		// 4. Save page data in DB
		c.db.StorePlaceholder(data.links)

	}
}

func (c *Crawler) crawlPage(url *url.URL) (pageData, error) {

	// 1. Get robots.txt
	resp, err := http.Get(url.Scheme + "://" + url.Host + "/robots.txt")
	if err != nil {
		return pageData{}, err
		// TODO: skip robots.txt stuff it one doesn't exist
	}
	defer resp.Body.Close()
	robotsTxt, err := io.ReadAll(resp.Body)
	if err != nil {
		return pageData{}, err
	}

	// 2. Check we can visit the page
	if !checkPageAllowed(string(robotsTxt), url.String()) {
		return pageData{}, fmt.Errorf("robots.txt forbids visiting page: %s", url)
	}

	// 3. Visit page
	resp, err = http.Get(url.String())
	if err != nil {
		return pageData{}, err
	}
	defer resp.Body.Close()

	//  4. Parse page contents
	data, err := parsePage(resp.Body, url)
	if err != nil {
		return pageData{}, err
	}

	return data, nil
}

func checkPageAllowed(robotsTxt, url string) bool {
	return grobotstxt.AgentAllowed(robotsTxt, crawlerName, url)
}
