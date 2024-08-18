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
	"github.com/z-riley/turdsearch/common/store"
	"github.com/z-riley/turdsearch/crawler/frontier"
	"github.com/z-riley/turdsearch/crawler/parser"
)

const (
	crawlerName = "TurdSeeker" // name of user agent in HTTP headers
)

var errDepthReached = errors.New("crawl depth reached")

type Crawler struct {
	frontier    *frontier.Frontier // queue of links to visit next
	db          *store.Storage     // database abstraction
	maxDepth    int                // maximum depth to recursively crawl
	gracePeriod time.Duration      // grace period before a webpage can get crawled again
}

func NewCrawler(gracePeriod time.Duration) (*Crawler, error) {

	db, err := store.NewStorageConn(store.StorageConfig{
		DatabaseName:          store.DatabaseName,
		CrawledDataCollection: store.CrawledDataCollection,
		WebgraphCollection:    store.WebgraphCollection,
		WordIndexCollection:   store.WordIndexCollection,
	})
	if err != nil {
		log.Fatal().Err(err)
	}

	c := &Crawler{
		frontier:    frontier.NewFrontier(),
		db:          db,
		maxDepth:    0,
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
		if err := c.frontier.Push(frontier.Link{URL: url, Depth: 1}); err != nil {
			return err
		}
	}
	return nil
}

// CrawlToTheDepth makes a crawler crawl until the desired depth is reached, at
// which point it will return
func (c *Crawler) CrawlToTheDepth(depth int) {
	c.maxDepth = depth
	for {
		err := c.crawlingSequence()
		if errors.Is(err, errDepthReached) {
			log.Warn().Err(err).Msg("Crawl depth reached")
			return
		} else if err != nil {
			log.Warn().Err(err).Msg("Crawl failed")
		}
	}
}

// crawlingSequence crawls a page from the frontier
func (c *Crawler) crawlingSequence() error {
	// 1. Get next link from frontier
	link, err := c.frontier.Dequeue()
	if err != nil {
		log.Error().Err(err).Msg("Failed to dequeue the next link from frontier")
		time.Sleep(1 * time.Second) // probably failed because queue was empty, so wait a second
		return err
	}

	// 2. Crawl page
	log.Info().Int("Depth", link.Depth).Msgf("Crawling page: %s", link.URL)
	url, err := url.Parse(link.URL)
	if err != nil {
		return err
	}
	data, err := c.crawlPage(url)
	if err != nil {
		return err
	}
	depth := link.Depth + 1
	log.Debug().Str("page", url.String()).Msgf("Found %d links", len(data.Links))

	// 3. Save page data to DB
	err = c.db.SavePageData(data)
	if err != nil {
		return err
	}

	// 4. Check if crawl depth reached
	if link.Depth >= c.maxDepth {
		return errDepthReached
	}

	// 5. Push new links into frontier
	for _, newURL := range data.Links {
		// Check if link was crawled recently
		isCrawledRecently, err := c.db.PageIsRecentlyCrawled(newURL, c.gracePeriod)
		if err != nil {
			return err
		}
		if !isCrawledRecently {
			if err := c.frontier.Push(frontier.Link{URL: newURL, Depth: depth}); err != nil {
				log.Error().Err(err)
			}
		} else {
			// log.Debug().Msgf("Not adding %s to frontier because it was already crawled in the last %v", url, c.gracePeriod)
		}
	}
	log.Info().Msgf("%d links added to frontier. New length: %d", len(data.Links), c.frontier.Len())

	// 6. Log frontier diagnostics
	m, err := c.frontier.TopNWebsites(10)
	if err != nil {
		return err
	}
	log.Info().Any("Top 10 websites", m).Msg("Frontier diagnostics")

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
	contentTypeHeader, ok := resp.Header["Content-Type"]
	if !ok {
		return store.PageData{}, fmt.Errorf("Not parsing page %s because of Content-Type header missing", url.String())
	}
	contentType := contentTypeHeader[0]

	parsable := strings.Contains(contentType, "text/html") || strings.Contains(contentType, "text/plain")
	if !parsable {
		return store.PageData{}, fmt.Errorf("Not parsing page %s because of non-text content type: %s", url.String(), contentType)
	}
	data, err := parser.ParsePage(resp.Body, url, timeAccessed)
	if err != nil {
		return store.PageData{}, err
	}

	return data, nil
}

func checkPageAllowed(robotsTxt, url string) bool {
	return grobotstxt.AgentAllowed(robotsTxt, crawlerName, url)
}

func (c *Crawler) Destroy() {
	c.db.Destroy()
}
