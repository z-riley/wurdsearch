package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/common/logging"
	"github.com/zac460/turdsearch/search"
)

func main() {
	logging.SetUpLogger(false)

	handler, err := search.NewServer()
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	mux := http.NewServeMux()
	mux.Handle("GET /search/{query}", handler)
	http.ListenAndServe("0.0.0.0:8080", mux)
}
