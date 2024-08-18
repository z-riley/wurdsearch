package main

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/z-riley/turdsearch/common/logging"
	"github.com/z-riley/turdsearch/common/store"
	"github.com/z-riley/turdsearch/indexer"
)

func main() {
	logging.SetUpLogger(false)
	db, err := store.NewStorageConn(store.StorageConfig{
		DatabaseName:          store.DatabaseName,
		CrawledDataCollection: store.CrawledDataCollection,
		WebgraphCollection:    store.WebgraphCollection,
		WordIndexCollection:   store.WordIndexCollection,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to DB")
	}
	defer db.Destroy()

	// Word index
	log.Info().Msg("Generating word index...")
	start := time.Now()
	w, err := indexer.NewWordIndexer(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to make new word indexer")
	}
	if err := w.GenerateWordIndex(db.Config.CrawledDataCollection); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate word index")
	}
	log.Info().Msgf("Generated word index in %v", time.Since(start))

	// Web graph
	log.Info().Msg("Generating web graph...")
	start = time.Now()
	g := indexer.NewWebgrapher(db)
	if err := g.GenerateWebgraph(); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate word index")
	}
	log.Info().Msgf("Generated web graph in %v", time.Since(start))

}
