package main

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	// "github.com/zac460/herolog"
)

var log zerolog.Logger

func main() {

	c, err := newCrawler()
	if err != nil {
		log.Fatal().Err(err)
	}
	defer c.destroy()

	seeds := []string{
		"https://google.com/",
		"https://reddit.com/",
		"https://example.com/",
	}
	if err := c.setSeeds(seeds); err != nil {
		log.Fatal().Err(err)
	}

	if err := c.crawlForever(); err != nil {
		log.Fatal().Err(err)
	}

}

func UNUSED(x ...any) {}

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

	zerolog.TimeFieldFormat = time.RFC3339Nano
	log = zerolog.New(multiWriter).With().Timestamp().Logger()
}
