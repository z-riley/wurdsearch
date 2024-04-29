package search

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/zac460/turdsearch/store"
)

func TestTheta(t *testing.T) {
	a := vector{
		label: "a",
		val: map[string]float64{
			"one":   1,
			"two":   2,
			"three": 3,
		},
	}
	b := vector{
		label: "b",
		val: map[string]float64{
			"one":   1,
			"two":   2,
			"three": 3,
		},
	}
	th := theta(a, b)
	expected := 0.0
	if th != expected {
		t.Errorf("Actual (%f) did not equal expected (%f)", th, expected)
	}

	a = vector{
		label: "a",
		val: map[string]float64{
			"one":   1,
			"two":   2,
			"three": 3,
		},
	}
	b = vector{
		label: "b",
		val: map[string]float64{
			"three": -3,
			"one":   -1,
			"two":   -2,
		},
	}
	th = theta(a, b)
	expected = math.Pi
	if th != expected {
		t.Errorf("Actual (%f) did not equal expected (%f)", th, expected)
	}
}

func TestDotProduct(t *testing.T) {
	a := vector{
		label: "a",
		val: map[string]float64{
			"one":   1,
			"two":   2,
			"three": 3,
			"four":  4,
		},
	}
	b := vector{
		label: "b",
		val: map[string]float64{
			"four":  4,
			"two":   2,
			"three": 3,
			"one":   1,
		},
	}

	dp := dotProduct(a, b)
	expected := 30.0
	if dp != expected {
		t.Error("Actual did not equal expected")
	}
}

func TestMod(t *testing.T) {
	v := vector{
		label: "a",
		val: map[string]float64{
			"one":   1,
			"two":   2,
			"three": 3,
			"four":  -4,
		},
	}
	mod := v.mod()
	expected := math.Sqrt(30.0)
	if mod != expected {
		t.Error("Actual did not equal expected")
	}

	v = vector{
		label: "a",
		val: map[string]float64{
			"one": 0,
			"two": 0,
		},
	}
	mod = v.mod()
	expected = 0.0
	if mod != expected {
		t.Error("Actual did not equal expected")
	}
}

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
