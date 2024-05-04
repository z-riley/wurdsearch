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

	wordIndexer, err := NewWordIndexer(db)
	if err != nil {
		return nil, err
	}

	return &Indexer{
		webgrapher:  NewWebgrapher(db),
		wordIndexer: wordIndexer,
		db:          db,
	}, nil
}
