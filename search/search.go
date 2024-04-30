package search

import (
	"errors"
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
func mergeScores(scores []pageScores, weights []float64) (pageScores, error) {
	if len(scores) != len(weights) {
		return pageScores{}, errors.New("pageScores and weights length mismatch")
	}

	// Extract every URL from page scores
	// Exploit that maps contain unique keys to only get one of each URL
	allUrls := make(map[string]float64)
	for _, score := range scores {
		for url := range score {
			allUrls[url] = 0.0
		}
	}

	// Add up scores, accounting for weights
	output := make(pageScores)
	for url := range allUrls {
		for i, score := range scores {
			val, exists := score[url]
			if exists {
				output[url] += (val * weights[i])
			}
		}
	}
	return output, nil
}
