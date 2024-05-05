package search

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestTFIDF(t *testing.T) {
	s, err := NewSearcher(getTestDB(t))
	if err != nil {
		t.Fatal(err)
	}

	query := "I have varroa"
	results, err := s.TFIDF(query)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(results)
}
func TestGenerateVector(t *testing.T) {

	start := time.Now()

	s, err := NewSearcher(getTestDB(t))
	if err != nil {
		t.Fatal(err)
	}

	searchTerm := []string{"varroa", "advice"}

	// First calculate IDF for each word
	IDFs := make(map[string]float64)
	for _, word := range searchTerm {
		IDF, err := s.inverseDocumentFrequency(word)
		if err != nil {
			t.Fatal(err)
		}
		IDFs[word] = IDF
	}

	url := "https://www.varroaresistant.uk/advice"
	result, err := s.generateVector(url, searchTerm, IDFs)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	url = "https://en.wikipedia.org/wiki/Wikipedia:File_Upload_Wizard"
	result, err = s.generateVector(url, searchTerm, IDFs)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	fmt.Println("Time:", time.Since(start))
}

func TestSearchTermVector(t *testing.T) {
	s, err := NewSearcher(getTestDB(t))
	if err != nil {
		t.Fatal(err)
	}

	searchWords := strings.Split("the quick quick fox", " ")
	result, err := s.queryVector(searchWords)
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
