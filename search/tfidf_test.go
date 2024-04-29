package search

import (
	"fmt"
	"strings"
	"testing"

	"github.com/zac460/turdsearch/store"
)

func TestTFIDF(t *testing.T) {
	s, err := NewSearcher(getTestDB(t))
	if err != nil {
		t.Fatal(err)
	}

	err = s.TFIDF("good advice")
	if err != nil {
		t.Fatal(err)
	}

}
func TestGenerateVector(t *testing.T) {
	s, err := NewSearcher(getTestDB(t))
	if err != nil {
		t.Fatal(err)
	}

	url := "https://www.varroaresistant.uk/advice"
	searchTerm := []string{"varroa advice"}
	result, err := s.generateVector(url, searchTerm)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	url = "https://en.wikipedia.org/wiki/Wikipedia:File_Upload_Wizard"
	searchTerm = []string{"varroa advice"}
	result, err = s.generateVector(url, searchTerm)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestSearchTermVector(t *testing.T) {
	s, err := NewSearcher(getTestDB(t))
	if err != nil {
		t.Fatal(err)
	}

	searchWords := strings.Split("the quick quick fox", " ")
	result, err := s.searchTermVector(searchWords)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(result)
}

func TestGetEveryRelevantDoc(t *testing.T) {
	s, err := NewSearcher(getTestDB(t))
	if err != nil {
		t.Fatal(err)
	}

	URLs, err := s.getEveryRelevantDoc([]string{"usually", "started"})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(URLs)

}
func TestTermFrequencies(t *testing.T) {
	s, err := NewSearcher(getTestDB(t))
	if err != nil {
		t.Fatal(err)
	}

	result, err := s.termFrequencies("local")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestInverseDocumentFrequency(t *testing.T) {
	s, err := NewSearcher(getTestDB(t))
	if err != nil {
		t.Fatal(err)
	}

	result, err := s.inverseDocumentFrequency("local")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

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
