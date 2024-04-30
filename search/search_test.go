package search

import (
	"reflect"
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

func TestMergeScores(t *testing.T) {
	a := pageScores{
		"url1": 1,
		"url2": 2,
		"url3": 3,
		"url4": 4,
	}
	b := pageScores{
		"url3": 1,
		"url4": 2,
		"url5": 3,
		"url6": 4,
	}
	weights := []float64{1.0, 4.0}
	result, err := mergeScores([]pageScores{a, b}, weights)
	if err != nil {
		t.Fatal(err)
	}
	expected := pageScores{
		"url1": 1,
		"url2": 2,
		"url3": 7,
		"url4": 12,
		"url5": 12,
		"url6": 16,
	}
	if !reflect.DeepEqual(result, expected) {
		t.Error("Result not equal to expected")
	}

	c := pageScores{
		"url1": 2,
		"url2": 4,
		"url6": 6,
		"url7": 8,
	}
	weights = []float64{1.0, 4.0, 0.5}
	result, err = mergeScores([]pageScores{a, b, c}, weights)
	if err != nil {
		t.Fatal(err)
	}
	expected = pageScores{
		"url1": 2,
		"url2": 4,
		"url3": 7,
		"url4": 12,
		"url5": 12,
		"url6": 19,
		"url7": 4,
	}
	if !reflect.DeepEqual(result, expected) {
		t.Error("Result not equal to expected")
	}
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
