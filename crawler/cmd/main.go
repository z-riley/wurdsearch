package main

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	// "github.com/zac460/herolog"
	"github.com/zac460/turdsearch/crawler"
)

const (
	crawlGracePeriod = 10 * time.Second
)

func main() {
	setUpLogger(false)
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

	c.CrawlForever()
}

// TODO: this should be in a common package
func setUpLogger(httpLogging bool) {
	var multiWriter io.Writer

	if httpLogging {
		multiWriter = io.MultiWriter(
			zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
			// herolog.NewLogHTTPWriter("http://0.0.0.0:2021", true),
		)
	} else {
		multiWriter = io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	log.Logger = zerolog.New(multiWriter).With().Timestamp().Caller().Logger()
}
