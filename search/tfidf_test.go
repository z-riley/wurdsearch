package search

import (
	"fmt"
	"testing"

	"github.com/zac460/turdsearch/store"
)

func TestWordFrequency(t *testing.T) {
	db, err := store.NewStorageConn(store.StorageConfig{
		DatabaseName:          store.DatabaseName,
		CrawledDataCollection: store.CrawledDataTestCollection,
		WebgraphCollection:    store.WebgraphTestCollection,
		WordIndexCollection:   store.WordIndexTestCollection,
	})
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewSearcher(db)
	if err != nil {
		t.Fatal(err)
	}

	result, err := s.wordFrequency("local")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestInverseDocumentFrequency(t *testing.T) {}
