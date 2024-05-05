package main

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/common/logging"
	"github.com/zac460/turdsearch/crawler"
)

const (
	parallelCrawlers = 5 // 5 is good, 10 isn't much better
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
	for n := 0; n < parallelCrawlers; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.CrawlForever()
		}()
	}
	wg.Wait()
	panic("All crawler Goroutines crashed somehow")
}
