package main

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/common/store"
	"github.com/zac460/turdsearch/indexer"
)

func main() {
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

	w, err := indexer.NewWordIndexer(db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to make new word indexer")
	}

	start := time.Now()

	if err := w.GenerateWordIndex(db.Config.CrawledDataCollection); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate word index")
	}

	log.Info().Msgf("Generated word index in %v", time.Since(start))
}
