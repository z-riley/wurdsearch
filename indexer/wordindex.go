package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zac460/turdsearch/store"
)

type WordIndexer struct {
	db *store.Storage
}

func NewWordIndexer(db *store.Storage) *WordIndexer {
	return &WordIndexer{
		db: db,
	}
}

// GenerateWordIndex generates a word index from crawled page data
func (w *WordIndexer) GenerateWordIndex() error {

	// Iterate through all crawled URLS
	if err := w.db.InitIterator(store.CrawledDataCollection); err != nil {
		return err
	}
	for {
		pageData, more, err := w.db.NextPageData()
		if err != nil {
			return err
		}
		if !more {
			break
		}

		// Update word index for each word on the page
		// This is horribly slow and should be improved at some point
		wordMap := make(map[string]uint)
		for _, word := range sanitiseString(strings.ToLower(pageData.Content)) {
			wordMap[word] += 1
		}
		for word, count := range wordMap {
			// Upsert the word in the DB with new {page url, how many times the word appears on that page}
			if err := w.db.UpdateWordReference(word, pageData.Url, count); err != nil {
				return fmt.Errorf("Update word index for %s: %v", word, err)
			}
		}

	}

	return nil
}

// sanitiseString removes grammar and switches to lowercase. Strings returned separately as slice
func sanitiseString(input string) []string {
	re := regexp.MustCompile("[a-zA-Z0-9'-]+")
	matches := re.FindAllString(input, -1)
	return matches
}
