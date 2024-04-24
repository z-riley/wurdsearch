package crawler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jimsmart/grobotstxt"
	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/crawler/frontier"
	"github.com/zac460/turdsearch/crawler/parser"
	"github.com/zac460/turdsearch/store"
)

const (
	crawlerName = "TurdSeeker" // name of user agent in HTTP headers
)

type Crawler struct {
	frontier    *frontier.Frontier // queue of links to visit next
	db          *store.Storage     // database abstraction
	gracePeriod time.Duration      // grace period before a webpage can get crawled again
}

func NewCrawler(gracePeriod time.Duration) (*Crawler, error) {

	db, err := store.NewStorageConn(store.StorageConfig{
		DatabaseName:          "turdsearch",
		CrawledDataCollection: "crawled_data",
		IndexedDataCollection: "indexed_data",
	})
	if err != nil {
		log.Fatal().Err(err)
	}

	c := &Crawler{
		frontier:    frontier.NewFrontier(),
		db:          db,
		gracePeriod: gracePeriod,
	}

	return c, nil
}

// SetSeeds takes URLs to use as starting points for crawling
func (c *Crawler) SetSeeds(urls []string) error {
	if len(urls) == 0 {
		return errors.New("No seed URLs provided")
	}

	for _, url := range urls {
		if err := c.frontier.Push(url); err != nil {
			return err
		}
	}
	return nil
}

func (c *Crawler) CrawlForever() {
	for {
		if err := c.crawlingSequence(); err != nil {
			log.Warn().Err(err).Msg("Crawl failed")
		}
	}
}

func (c *Crawler) crawlingSequence() error {
	// 1. Get next link from frontier
	link, err := c.frontier.Dequeue()
	if err != nil {
		log.Error().Err(err).Msg("Failed to dequeue the next link from frontier")
		return err
	}
	link, ok := link.(string)
	if !ok {
		return fmt.Errorf("%v is not castable to string", link)
	}

	// 2. Crawl page
	log.Info().Msgf("Crawling page: %s", link.(string))
	url, err := url.Parse(link.(string))
	if err != nil {
		return err
	}
	data, err := c.crawlPage(url)
	if err != nil {
		return err
	}
	log.Debug().Str("page", url.String()).Msgf("Found %d links", len(data.Links))

	// 3. Put new links into frontier if the pages haven't been scraped recently
	for _, link := range data.Links {

		// Check if link was crawled recently
		isCrawledRecently, err := c.db.PageIsRecentlyCrawled(link, c.gracePeriod)
		if err != nil {
			return err
		}
		if !isCrawledRecently {
			if err := c.frontier.Push(link); err != nil {
				log.Error().Err(err)
			}
		} else {
			log.Debug().Msgf("Not adding %s to frontier because it was already crawled in the last %v", url, c.gracePeriod)
		}
	}
	log.Info().Msgf("%d links added to frontier. New length: %d", len(data.Links), c.frontier.Len())

	// 4. Save page data to DB
	err = c.db.SavePageData(data)
	if err != nil {
		return err
	}

	// (5. Log frontier diagnostics)
	q, err := c.frontier.GetAll()
	if err != nil {
		return err
	}
	m, err := frontier.CountOccurrances(q)
	if err != nil {
		return err
	}
	log.Info().Any("map", m).Msg("Frontier diagnostics")

	return nil
}

// crawlPage crawls a page, obeying robots.txt
func (c *Crawler) crawlPage(url *url.URL) (store.PageData, error) {

	// 1. Get robots.txt
	resp, err := http.Get(url.Scheme + "://" + url.Host + "/robots.txt")
	if err != nil {
		return store.PageData{}, err
	}
	if resp.StatusCode == 404 {
		log.Debug().Msgf("robots.txt not found for %s. Continuing anyway", url.Host)
	}
	defer resp.Body.Close()
	robotsTxt, err := io.ReadAll(resp.Body)
	if err != nil {
		return store.PageData{}, err
	}

	// 2. Check we can visit the page
	if !checkPageAllowed(string(robotsTxt), url.String()) {
		return store.PageData{}, fmt.Errorf("robots.txt forbids visiting page: %s", url)
	}

	// 3. Visit page
	timeAccessed := time.Now()
	resp, err = http.Get(url.String())
	if err != nil {
		return store.PageData{}, err
	}
	defer resp.Body.Close()

	//  4. Parse page contents
	contentType := resp.Header["Content-Type"][0]
	parsable := strings.Contains(contentType, "text/html") || strings.Contains(contentType, "text/plain")
	if !parsable {
		return store.PageData{}, fmt.Errorf("Not parsing page %s because of non-text content type: %s", url.String(), contentType)
	}
	data := parser.ParsePage(resp.Body, url, timeAccessed)

	return data, nil
}

func checkPageAllowed(robotsTxt, url string) bool {
	return grobotstxt.AgentAllowed(robotsTxt, crawlerName, url)
}

func (c *Crawler) Destroy() {
	c.db.Destroy()
}
