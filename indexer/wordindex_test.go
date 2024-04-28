package main

import (
	"testing"

	"github.com/zac460/turdsearch/store"
)

func TestGenerateWordIndex(t *testing.T) {
	db, err := store.NewStorageConn(store.StorageConfig{
		DatabaseName:          store.DatabaseName,
		CrawledDataCollection: store.CrawledDataTestCollection,
		WebgraphCollection:    store.WebgraphTestCollection,
		WordIndexCollection:   store.WordIndexTestCollection,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()

	w := NewWordIndexer(db)

	if err := w.GenerateWordIndex(db.Config.CrawledDataCollection); err != nil {
		t.Fatal(err)
	}
}
