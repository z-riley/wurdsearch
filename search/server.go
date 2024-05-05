package search

import (
	"encoding/json"
	"net/http"

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
		WordIndexCollection:   store.WordIndexTestCollection,
	})
	if err != nil {
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

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	query := r.PathValue("query")

	result, err := s.searcher.Search(query)
	if err != nil {
		http.Error(w, "Search failed: "+err.Error(), http.StatusInternalServerError)
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(bytes)
}
