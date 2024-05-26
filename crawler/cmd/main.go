package main

import (
	"slices"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/common/logging"
	"github.com/zac460/turdsearch/crawler"
)

const (
	/* Crawling test 60s:
	5 crawlers: 500
	10 crawlers: 734 (CPU 85%)
	20 crawlers: 869 (CPU 100%)
	*/
	parallelCrawlers = 10
	crawlDepth       = 3
	crawlGracePeriod = 10 * time.Second
)

func main() {
	logging.SetUpLogger(false)
	log.Info().Msg("Begin")

	c, err := crawler.NewCrawler(crawlGracePeriod)
	if err != nil {
		log.Fatal().Err(err)
	}
	defer c.Destroy()

	seeds := []string{
		"https://en.wikipedia.org/wiki/Imperial_Crypt",
		"https://puginarug.com/",
		"https://www.varroaresistant.uk/",
		"https://parade.com/1116816/marynliles/fun-websites/",
	}
	if err := c.SetSeeds(seeds); err != nil {
		log.Fatal().Err(err).Msg("Failed to set seeds")
	}

	var wg sync.WaitGroup
	crawlers := slices.Min([]int{parallelCrawlers, len(seeds)})
	for n := 0; n < crawlers; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.CrawlToTheDepth(crawlDepth)
		}()
	}
	wg.Wait()
	log.Info().Msg("Crawl ended")
}
