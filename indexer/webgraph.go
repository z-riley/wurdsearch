package main

import (
	"fmt"

	"github.com/zac460/turdsearch/common/store"
)

type Webgrapher struct {
	db *store.Storage
}

func NewWebgrapher(db *store.Storage) *Webgrapher {
	return &Webgrapher{
		db: db,
	}
}

// GenerateWebgraph generates a webgraph from the crawled data in the database
func (w *Webgrapher) GenerateWebgraph() error {

	if err := w.db.InitIterator(w.db.Config.CrawledDataCollection); err != nil {
		return err
	}

	// Iterate over every URL in the crawled data collection
	for {
		data, more, err := w.db.NextPageData()
		if err != nil {
			return err
		}
		if !more {
			break
		}

		// Populate linksTo field
		if err := w.db.SaveNode(store.Node{
			Url:       data.Url,
			LinksTo:   data.Links,
			LinksFrom: []string{"todo"},
		}); err != nil {
			return fmt.Errorf("Failed to save node: %v", err)
		}

		// Note: not possible to populate linksFrom field because it is N^2 database
		// calls, where N is the number of crawled documents. Leaving as todo for now

	}
	return nil
}
