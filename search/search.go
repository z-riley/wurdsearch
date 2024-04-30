package search

import (
	"fmt"
	"strings"

	"github.com/zac460/turdsearch/store"
)

type Searcher struct {
	db *store.Storage
}

func NewSearcher(store *store.Storage) (*Searcher, error) {
	return &Searcher{
		db: store,
	}, nil
}

type Result struct {
	Url               string  `json:"url"`
	ConfidencePercent float64 `json:"confidence"`
}

// Search executes a search, returning a slice of relevant documents
func (s *Searcher) Search(query string) ([]store.PageData, error) {

	query = sanitiseQuery(query)

	// TF-IDF
	results, err := s.TFIDF(query)
	if err != nil {
		return []store.PageData{}, err
	}

	for url, score := range results {
		fmt.Printf("%2.f%%\t%s\n", score, url)
	}

	// Do weighted sum with other search algorithms once they're implemented

	return []store.PageData{}, nil
}

// sanitiseQuery sanitise a search query before use in search algorithms
func sanitiseQuery(query string) string {
	query = strings.ToLower(query)
	query = strings.TrimSpace(query)
	return query
}

// pageScores holds number scores for URLs
type pageScores map[string]float64

// mergeScores combines serach results according to their given weightings
func mergeScores(...pageScores) {

}
