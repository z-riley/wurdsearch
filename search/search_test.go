package search

import (
	"testing"

	"github.com/zac460/turdsearch/store"
)

func TestSearch(t *testing.T) {
	s, err := NewSearcher(getTestDB(t))
	if err != nil {
		t.Fatal(err)
	}

	query := "bees varroa"
	results, err := s.Search(query)
	if err != nil {
		t.Fatal(err)
	}
	UNUSED(results)
}

func getTestDB(t *testing.T) *store.Storage {
	db, err := store.NewStorageConn(store.StorageConfig{
		DatabaseName:          store.DatabaseName,
		CrawledDataCollection: store.CrawledDataTestCollection,
		WebgraphCollection:    store.WebgraphTestCollection,
		WordIndexCollection:   store.WordIndexTestCollection,
	})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func UNUSED(x ...any) {}
