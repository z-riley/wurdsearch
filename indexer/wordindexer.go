package indexer

import (
	"regexp"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/common/lemmatiser"
	"github.com/zac460/turdsearch/common/stopwords"
	"github.com/zac460/turdsearch/common/store"
)

const indexerWorkers = 10

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

	// Set up semaphore
	var wg sync.WaitGroup
	sem := make(chan struct{}, indexerWorkers)

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

		wg.Add(1)
		sem <- struct{}{} // take from semaphore

		go func(pageData store.PageData) {
			defer wg.Done()
			defer func() { <-sem }() // give to semphore

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
					log.Error().Msgf("Failed to update word references for %s: %v", pageData.Url, err)
				}
			}
		}(pageData)
	}

	wg.Wait()

	return nil
}

// sanitiseString removes grammar and switches to lowercase. Strings returned separately as slice
func sanitiseString(input string) []string {
	re := regexp.MustCompile("[a-zA-Z0-9'-]+")
	matches := re.FindAllString(input, -1)
	return matches
}
