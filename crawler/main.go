package main

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	// "github.com/zac460/herolog"
)

var log zerolog.Logger

func main() {
	setUpLogger(false)

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
		log.Fatal().Err(err).Msg("Failed to set seeds")
	}

	c.crawlForever()
}

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

	log = zerolog.New(multiWriter).With().Timestamp().Caller().Logger()
}

func UNUSED(x ...any) {}
