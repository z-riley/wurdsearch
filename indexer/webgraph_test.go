package main

import (
	"testing"

	"github.com/zac460/turdsearch/store"
)

func TestGenerateWebgraph(t *testing.T) {
	db, err := store.NewStorageConn(store.StorageConfig{
		DatabaseName:          store.DatabaseName,
		CrawledDataCollection: store.CrawledDataTestCollection,
		WebgraphCollection:    store.WebgraphTestCollection,
	})
	if err != nil {
		t.Error(err)
	}
	defer db.Destroy()

	w := NewWebgrapher(db)

	err = w.GenerateWebgraph()
	if err != nil {
		t.Error(err)
	}
}
