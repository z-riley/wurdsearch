package indexer

import (
	"testing"

	"github.com/z-riley/wurdsearch/common/store"
)

func TestGenerateWordIndex(t *testing.T) {
	db, err := store.NewStorageConn(store.Config{
		DatabaseName:          store.DatabaseName,
		CrawledDataCollection: store.CrawledDataCollection,
		WebgraphCollection:    store.WebgraphTestCollection,
		WordIndexCollection:   store.WordIndexTestCollection,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()

	w, err := NewWordIndexer(db)
	if err != nil {
		t.Fatal(err)
	}

	if err := w.GenerateWordIndex(db.Config.CrawledDataCollection); err != nil {
		t.Fatal(err)
	}
}
