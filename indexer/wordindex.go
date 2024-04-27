package main

import "github.com/zac460/turdsearch/store"

// WordEntry contains which websites use a particular word, and how many times
// it appears on each page
type WordEntry struct {
	word       string          `bson:"word"`
	references map[string]uint `bson:"references"`
}

type WordIndexer struct {
	db *store.Storage
}

func NewWordIndexer(db *store.Storage) *WordIndexer {
	return &WordIndexer{
		db: db,
	}
}
