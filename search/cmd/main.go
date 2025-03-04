package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/z-riley/wurdsearch/common/logging"
	"github.com/z-riley/wurdsearch/search"
)

const port = 8080

func main() {
	logging.Init()

	handler, err := search.NewServer()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	mux := http.NewServeMux()
	mux.Handle("GET /search/{query}", handler)
	log.Info().Msgf("Starting server on port %d", port)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", 8080), mux)
}
