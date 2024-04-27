package main

import (
	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/store"
)

type Indexer struct {
	webgrapher  *Webgrapher
	wordIndexer *WordIndexer
	db          *store.Storage
}

func NewIndexer() (*Indexer, error) {

	db, err := store.NewStorageConn(store.StorageConfig{
		DatabaseName:          store.DatabaseName,
		CrawledDataCollection: store.CrawledDataCollection,
		WebgraphCollection:    store.WebgraphCollection,
		WordIndexCollection:   store.WordIndexCollection,
	})
	if err != nil {
		log.Fatal().Err(err)
	}

	webgraph := NewWebgrapher(db)
	wordIndexer := NewWordIndexer(db)

	return &Indexer{
		webgrapher:  webgraph,
		wordIndexer: wordIndexer,
		db:          db,
	}, nil
}
