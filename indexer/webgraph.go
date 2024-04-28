package main

import (
	"fmt"

	"github.com/zac460/turdsearch/store"
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

	if err := w.db.InitIterator(store.CrawledDataCollection); err != nil {
		return err
	}

	// Iterate over every URL in the crawled data collection
	for {
		data, more, err := w.db.IterateNext()
		if err != nil {
			return err
		}
		if !more {
			break
		}

		// TODO: Populate linksFrom field (any way to do this not N^2?)

		if err := w.db.SaveNode(store.Node{
			Url:       data.Url,
			LinksTo:   data.Links,
			LinksFrom: []string{"todo"},
		}); err != nil {
			return fmt.Errorf("Failed to save node: %v", err)
		}

	}
	return nil
}
