package indexer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/common/lemmatiser"
	"github.com/zac460/turdsearch/common/stopwords"
	"github.com/zac460/turdsearch/common/store"
)

type WordIndexer struct {
	lemmatiser *lemmatiser.Lemmatiser
	db         *store.Storage
}

func NewWordIndexer(db *store.Storage) (*WordIndexer, error) {
	l, err := lemmatiser.NewLemmatiser()
	if err != nil {
		return nil, err
	}

	return &WordIndexer{
		lemmatiser: l,
		db:         db,
	}, nil
}

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
		pageData, more, err := w.db.NextPageData()
		if err != nil {
			return err
		}
		if !more {
			break
		}

		count++
		log.Info().Msgf("Generating word index. Progress: %d/%d", count, length)

		// Update word index for each word on the page
		wordCounts := make(map[string]uint)
		words := sanitiseString(strings.ToLower(pageData.Content))
		for _, word := range words {
			// Ignore stop words since they'll never be searched for anyway
			if stopwords.IsStopWord(word) {
				continue
			}
			// Only store lemmas
			lemmatisedWord := w.lemmatiser.Lemmatise(word)
			wordCounts[lemmatisedWord] += 1
		}
		if len(wordCounts) > 0 {
			err = w.db.UpdateWordReferences(pageData.Url, wordCounts, uint(len(words)))
			if err != nil {
				return fmt.Errorf("Failed to update word references for %s: %v", pageData.Url, err)
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
