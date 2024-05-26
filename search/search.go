package search

import (
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/common/lemmatiser"
	"github.com/zac460/turdsearch/common/store"
)

type PageData struct {
	Score        float64   `bson:"score"`
	LastAccessed time.Time `bson:"lastAccessed"`
	Content      string    `bson:"content"`
}

// PageScores holds number scores for URLs
type PageScores map[string]float64

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
func (s *Searcher) Search(query string) (map[string]PageData, error) {
	start := time.Now()
	query = sanitiseQuery(query)

	// TF-IDF
	TFIDFScores, err := s.TFIDF(query)
	if err != nil {
		return nil, err
	}

	// Do weighted sum with other search algorithms once they're implemented
	finalScores, err := mergeScores(
		[]PageScores{TFIDFScores},
		[]float64{1.0},
	)
	if err != nil {
		return nil, err
	}

	log.Info().Msgf("Found %d results for '%s', in %dms", len(finalScores), query, time.Since(start).Milliseconds())

	// Get accompanying page data from DB
	results := make(map[string]PageData)
	for URL, score := range finalScores {
		pageData, err := s.db.FetchPageData(URL)
		if err != nil {
			return nil, err
		}
		const contentTrim = 150
		results[URL] = PageData{
			Score:        score,
			LastAccessed: pageData.LastAccessed,
			Content:      truncate(pageData.Content, 150),
		}
	}

	return results, nil
}

// sanitiseQuery sanitises a search query before use in search algorithms
func sanitiseQuery(query string) string {
	query = strings.ToLower(query)
	query = strings.TrimSpace(query)

	// Handle symbols
	result := ""
	for _, letter := range query {
		switch letter {
		case '&', '/', '\\', '-', '+', '=', '_':
			result += " "
		case '!', '?', '"', '^', '(', ')', '{', '}', '[', ']', '<', '>', ',', '.':
			result += ""
		default:
			result += string(letter)
		}
	}
	return result
}

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

// truncate takes the first n characters of a string
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return strings.TrimSpace(s[:n-3]) + "..."
}
