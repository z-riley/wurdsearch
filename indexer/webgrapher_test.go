package indexer

import (
	"testing"

	"github.com/z-riley/turdsearch/common/store"
)

func TestGenerateWebgraph(t *testing.T) {
	db, err := store.NewStorageConn(store.Config{
		DatabaseName:          store.DatabaseName,
		CrawledDataCollection: store.CrawledDataTestCollection,
		WebgraphCollection:    store.WebgraphTestCollection,
		WordIndexCollection:   store.WordIndexTestCollection,
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
