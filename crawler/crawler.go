package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/jimsmart/grobotstxt"
)

const crawlerName = "TurdSeeker"

type Crawler struct {
	frontier *frontier // queue of links to visit next
	db       *storage  // database abstraction
}

func newCrawler() (*Crawler, error) {

	db, err := newStorageConn()
	if err != nil {
		log.Fatal().Err(err)
	}

	c := &Crawler{
		frontier: newFrontier(),
		db:       db,
	}

	return c, nil
}

// setSeeds takes URLs to use as starting points for crawling
func (c *Crawler) setSeeds(urls []string) error {
	if len(urls) == 0 {
		return errors.New("No seed URLs provided")
	}

	for _, url := range urls {
		if err := c.frontier.push(url); err != nil {
			return err
		}
	}
	return nil
}

func (c *Crawler) crawlForever() {

	for {
		if err := c.crawlingSequence(); err != nil {
			log.Error().Err(err).Msg("Crawl failed")
		}
	}
}

func (c *Crawler) crawlingSequence() error {
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
	log.Debug().Msgf("Crawling page: %s", link.(string))
	url, err := url.Parse(link.(string))
	if err != nil {
		return err
	}
	data, err := c.crawlPage(url)
	if err != nil {
		return err
	}
	log.Debug().Str("page", url.String()).Msgf("Found %d links", len(data.Links))
	fmt.Printf("Data: %+v\n", data)

	// 3. Put new links into frontier
	for _, link := range data.Links {
		if err := c.frontier.push(link); err != nil {
			log.Error().Err(err)
		}
	}
	log.Debug().Msgf("Queue length: %v", c.frontier.queue.GetLen())

	// 4. Save page data to DB
	err = c.db.savePageData(data)
	if err != nil {
		return err
	}

	// (5. log frontier diagnostics)
	q, err := c.frontier.getAll()
	if err != nil {
		return err
	}
	m, err := countOccurrances(q)
	if err != nil {
		return err
	}
	log.Info().Any("map", m).Msg("Diagnostics")

	return nil
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
	timeAccessed := time.Now()
	resp, err = http.Get(url.String())
	if err != nil {
		return pageData{}, err
	}
	defer resp.Body.Close()

	//  4. Parse page contents
	data := parsePage(resp.Body, url, timeAccessed)

	return data, nil
}

func checkPageAllowed(robotsTxt, url string) bool {
	return grobotstxt.AgentAllowed(robotsTxt, crawlerName, url)
}

func (c *Crawler) destroy() {
	c.db.destroy()
}
