package main

import "github.com/zac460/turdsearch/store"

type WordIndexer struct {
	db *store.Storage
}

func NewWordIndexer(db *store.Storage) *WordIndexer {
	return &WordIndexer{
		db: db,
	}
}

func (w *WordIndexer) GenerateWordIndex() error {

	// 1. Iterate through all crawled URLS

	// 2. For each word on the page:
	// Upsert the word in the DB with new {page url, how many times the word appears on that page}

	return nil
}
