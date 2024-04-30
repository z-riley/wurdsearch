package search

import (
	"fmt"
	"reflect"
	"sort"
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

	// Print results in ascending order
	type Entry struct {
		Key   string
		Value float64
	}

	var entries []Entry
	for key, value := range results {
		entries = append(entries, Entry{Key: key, Value: value})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value < entries[j].Value // Descending order
	})
	for _, entry := range entries {
		fmt.Printf("%.3f%% - %s\n", entry.Value, entry.Key)
	}
}

func TestMergeScores(t *testing.T) {
	a := PageScores{
		"url1": 1,
		"url2": 2,
		"url3": 3,
		"url4": 4,
	}
	b := PageScores{
		"url3": 1,
		"url4": 2,
		"url5": 3,
		"url6": 4,
	}
	weights := []float64{1.0, 4.0}
	result, err := mergeScores([]PageScores{a, b}, weights)
	if err != nil {
		t.Fatal(err)
	}
	expected := PageScores{
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

	c := PageScores{
		"url1": 2,
		"url2": 4,
		"url6": 6,
		"url7": 8,
	}
	weights = []float64{1.0, 4.0, 0.5}
	result, err = mergeScores([]PageScores{a, b, c}, weights)
	if err != nil {
		t.Fatal(err)
	}
	expected = PageScores{
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
