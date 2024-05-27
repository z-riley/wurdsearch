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
	crawlDepth       = 5
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
		"https://www.cnn.com",
		"https://www.bbc.com/news",
		"https://www.nytimes.com",
		"https://www.reuters.com",
		"https://www.aljazeera.com",
		"https://techcrunch.com",
		"https://www.wired.com",
		"https://www.theverge.com",
		"https://www.cnet.com",
		"https://www.gizmodo.com",
		"https://www.scientificamerican.com",
		"https://www.nature.com",
		"https://www.nationalgeographic.com",
		"https://www.newscientist.com",
		"https://www.smithsonianmag.com",
		"https://www.webmd.com",
		"https://www.mayoclinic.org",
		"https://www.healthline.com",
		"https://www.medicalnewstoday.com",
		"https://www.cdc.gov",
		"https://www.khanacademy.org",
		"https://www.coursera.org",
		"https://www.edx.org",
		"https://ocw.mit.edu",
		"https://online.stanford.edu",
		"https://www.bloomberg.com",
		"https://www.marketwatch.com",
		"https://www.forbes.com",
		"https://www.wsj.com",
		"https://www.ft.com",
		"https://www.rollingstone.com",
		"https://www.ew.com",
		"https://www.tmz.com",
		"https://www.variety.com",
		"https://www.ign.com",
		"https://www.espn.com",
		"https://www.bleacherreport.com",
		"https://www.si.com",
		"https://www.foxsports.com",
		"https://www.bbc.com/sport",
		"https://www.lonelyplanet.com",
		"https://www.tripadvisor.com",
		"https://www.nationalgeographic.com/travel",
		"https://www.travelandleisure.com",
		"https://www.cntraveler.com",
		"https://www.goop.com",
		"https://www.thespruce.com",
		"https://www.realsimple.com",
		"https://www.lifehacker.com",
		"https://www.goodhousekeeping.com",
		"https://www.epicurious.com",
		"https://www.bonappetit.com",
		"https://www.seriouseats.com",
		"https://www.foodnetwork.com",
		"https://www.allrecipes.com",
		"https://www.vogue.com",
		"https://www.elle.com",
		"https://www.gq.com",
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
