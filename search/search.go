package search

import (
	"errors"
	"strings"

	"github.com/zac460/turdsearch/lemmatiser"
	"github.com/zac460/turdsearch/store"
)

type Searcher struct {
	lemmatiser *lemmatiser.Lemmatiser
	db         *store.Storage
}

func NewSearcher(store *store.Storage) (*Searcher, error) {
	l, err := lemmatiser.NewLemmatiser()
	if err != nil {
		return nil, err
	}

	return &Searcher{
		lemmatiser: l,
		db:         store,
	}, nil
}

// Search executes a search, returning a slice of relevant documents
func (s *Searcher) Search(query string) (PageScores, error) {

	query = sanitiseQuery(query)

	// TF-IDF
	TFIDFScores, err := s.TFIDF(query)
	if err != nil {
		return PageScores{}, err
	}

	// Do weighted sum with other search algorithms once they're implemented
	finalScores, err := mergeScores(
		[]PageScores{TFIDFScores},
		[]float64{1.0},
	)
	if err != nil {
		return PageScores{}, err
	}

	return finalScores, nil
}

// sanitiseQuery sanitise a search query before use in search algorithms
func sanitiseQuery(query string) string {
	query = strings.ToLower(query)
	query = strings.TrimSpace(query)
	return query
}

// PageScores holds number scores for URLs
type PageScores map[string]float64

// mergeScores combines serach results according to their given weightings
func mergeScores(scores []PageScores, weights []float64) (PageScores, error) {
	if len(scores) != len(weights) {
		return PageScores{}, errors.New("pageScores and weights length mismatch")
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
	output := make(PageScores)
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
