package search

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/common/store"
)

type Handler struct {
	searcher Searcher
}

func NewServer() (*Handler, error) {
	db, err := store.NewStorageConn(store.StorageConfig{
		DatabaseName:          store.DatabaseName,
		CrawledDataCollection: store.CrawledDataCollection,
		WebgraphCollection:    store.WebgraphCollection,
		WordIndexCollection:   store.WordIndexCollection,
	})
	if err != nil {
		log.Warn().Msg("Make sure DB is running")
		return nil, err
	}
	s, err := NewSearcher(db)
	if err != nil {
		return nil, err
	}

	return &Handler{
		searcher: *s,
	}, nil
}

// crap to be moved
type Text struct {
	Value  string `json:"value"`
	IsBold bool   `json:"is_bold"`
}
type Listing struct {
	Title   []Text `json:"title"`
	Extract []Text `json:"extract"`
	URL     string `json:"url"`
	Source  string `json:"source"`
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	query := r.PathValue("query")

	result, err := s.searcher.Search(query)
	if err != nil {
		http.Error(w, "Search failed: "+err.Error(), http.StatusInternalServerError)
	}

	// Assemble JSON for front-end
	var listings []Listing
	for URL, data := range result {
		listings = append(listings,
			Listing{
				Title:   []Text{{Value: data.Title, IsBold: false}},
				Extract: []Text{{Value: data.Content, IsBold: false}},
				URL:     URL,
				Source:  "Turdsearch",
			},
		)
	}

	bytes, err := json.Marshal(listings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(bytes)
}
