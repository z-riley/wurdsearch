package search

import "github.com/zac460/turdsearch/store"

type Searcher struct {
	db *store.Storage
}

func NewSearcher(store *store.Storage) (*Searcher, error) {
	return &Searcher{
		db: store,
	}, nil
}

// Search executes a search, returning a slice of relevant documents
func (s *Searcher) Search(term string) ([]store.PageData, error) {

	// Do weighted sum with other search algorithms once they're implemented

	return []store.PageData{}, nil
}
