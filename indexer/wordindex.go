package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
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
func (w *WordIndexer) GenerateWordIndex(collectionName string) error {

	// Iterate through all crawled URLS
	if err := w.db.InitIterator(collectionName); err != nil {
		return err
	}

	// Keep track of progress
	length, err := w.db.Len(collectionName)
	if err != nil {
		return err
	}
	count := 0

	for {
		count++
		log.Debug().Msgf("Generating word index. Progress: %d/%d", count, length)

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
		words := sanitiseString(strings.ToLower(pageData.Content))
		for _, word := range words {
			wordMap[word] += 1
		}
		for word, count := range wordMap {
			// Upsert the word in the DB with new data
			wordCount := uint(len(words))
			err := w.db.UpdateWordReference(word, pageData.Url, count, wordCount)
			if err != nil {
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
