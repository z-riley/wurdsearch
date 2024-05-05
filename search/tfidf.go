package search

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zac460/turdsearch/common/store"
	"go.mongodb.org/mongo-driver/mongo"
)

type VectorResult struct {
	Vec vector
	Err error
}

// TFIDF performs a TF-IDF search and returns relevant URLs and match confidence
// as a percentage
func (s *Searcher) TFIDF(query string) (PageScores, error) {

	// 1. Lemmatise and remove stop words
	var words []string
	rawWords := strings.Split(query, " ")
	for _, word := range rawWords {
		lword := s.lemmatiser.Lemmatise(word)
		if !isStopWord(lword) {
			words = append(words, lword)
		}
	}

	// 2. Calculate page vectors
	// a. Fetch all documents relevant to search query
	urls, err := s.getEveryRelevantDoc(words, 3)
	if err != nil {
		return PageScores{}, err
	}
	// b. Calculate IDF for each word
	IDFs := make(map[string]float64)
	for _, word := range words {
		IDF, err := s.inverseDocumentFrequency(word)
		if err != nil {
			return PageScores{}, err
		}
		IDFs[word] = IDF
	}
	// c. Calculate TF-IDF vectors
	start := time.Now() // temp
	vectors, err := s.generateAllVectors(urls, words, IDFs)
	if err != nil {
		return PageScores{}, fmt.Errorf("Failed to generate vectors: %v", err)
	}
	fmt.Println("Calculating all page vectors: ", time.Since(start)) // temp

	// 3. Get query vector
	queryVec, err := s.queryVector(words)
	if err != nil {
		return PageScores{}, err
	}
	log.Debug().Msgf("Seach vector: %v", queryVec)

	// 4. Compare the query vector to each page vector
	scores := make(PageScores)
	for _, pageVec := range vectors {
		theta := theta(queryVec, pageVec)
		percent := 100 * (1 - theta/math.Pi)
		scores[pageVec.label] = percent
	}

	return scores, nil
}

// generateAllVectors generates TF-IDF vectors for the words in the search query for each provided URL
func (s *Searcher) generateAllVectors(urls []string, queryWords []string, IDFs map[string]float64) ([]vector, error) {
	numWorkers := int(math.Floor(0.9 * float64(store.MaxConnections())))
	urlsChan := make(chan string, len(urls))
	resultsChan := make(chan VectorResult, len(urls))

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for n := 0; n < numWorkers; n++ {
		go func() {
			defer wg.Done()
			for url := range urlsChan {
				v, err := s.generateVector(url, queryWords, IDFs)
				resultsChan <- VectorResult{Vec: v, Err: err}
			}
		}()
	}
	for _, url := range urls {
		urlsChan <- url // send work to the worker pool
	}
	close(urlsChan) // indicate no more to be sent
	wg.Wait()
	close(resultsChan)

	var vectors []vector
	for result := range resultsChan {
		if result.Err != nil {
			return nil, result.Err
		}
		vectors = append(vectors, result.Vec)
	}
	return vectors, nil
}

// generateVector calcuates a vector for a given query on a page
func (s *Searcher) generateVector(url string, queryWords []string, wordIDFs map[string]float64) (vector, error) {
	v := vector{
		label: url,
		val:   make(map[string]float64),
	}

	for _, word := range queryWords {

		start := time.Now() // temp
		TFs, err := s.termFrequencies(word)
		if err != nil {
			return v, err
		}
		TF := TFs[url]
		fmt.Printf("Calculating term frequencies for %s on %s: %v\n", word, url, time.Since(start)) // temp

		IDF, ok := wordIDFs[word]
		if !ok {
			return v, fmt.Errorf("Supplied IDFs must include the word \"%s\"", word)
		}
		v.val[word] = TF * IDF
	}

	return v, nil
}

// queryVector gets the TF-IDF vector a query
func (s *Searcher) queryVector(words []string) (vector, error) {
	result := vector{
		label: "query",
		val:   map[string]float64{},
	}

	// Count occurances of each word
	wordCounts := make(map[string]uint)
	for _, word := range words {
		wordCounts[word] += 1
	}

	// Get TF-IDF of each word in search query
	for word, count := range wordCounts {
		TF := float64(count) / float64(len(words))
		IDF, err := s.inverseDocumentFrequency(word)
		if err != nil {
			return result, err
		}
		result.val[word] = TF * IDF
	}

	return result, nil
}

// getEveryRelevantDoc retrieves the URL of every document that contains n or more occurances
// of any of the words in the search query
func (s *Searcher) getEveryRelevantDoc(words []string, minOccurances uint) ([]string, error) {
	var urls []string
	for _, word := range words {
		wordIndexes, err := s.db.GetWordIndex(word)
		if errors.Is(err, mongo.ErrNoDocuments) {
			continue
		} else if err != nil {
			return urls, err
		}
		for url, ref := range wordIndexes.References {
			if ref.Count >= minOccurances {
				urls = append(urls, url)
			}
		}
	}
	return urls, nil
}

// termFrequencies returns the TFs of a specified word in every document.
// The TF is the number of times a word appears in a document divided by the total number
// of words in the document
func (s *Searcher) termFrequencies(word string) (PageScores, error) {
	TFs := make(PageScores)

	// Get all links which contain that word
	doc, err := s.db.GetWordIndex(word)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return TFs, nil
	} else if err != nil {
		return TFs, err
	}

	// Populate URL -> TF pairings for that word
	for url, ref := range doc.References {
		TF := float64(ref.Count) / float64(ref.Length)
		TFs[url] = TF
	}

	return TFs, nil
}

// inverseDocumentFrequency returns the IDF for a word. IDF is defined as the log of
// the number of documents divided by the number of documents that contain that word
func (s *Searcher) inverseDocumentFrequency(word string) (float64, error) {

	// Total number of documents
	numDocs, err := s.db.Len(s.db.Config.CrawledDataCollection)
	if err != nil {
		return 0.0, err
	}

	// Number of documents containing the word
	wordIndex, err := s.db.GetWordIndex(word)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return 0.0, nil
	} else if err != nil {
		return 0.0, err
	}
	numDocsContainingWord := len(wordIndex.References)

	return math.Log(float64(numDocs) / float64(numDocsContainingWord)), nil
}
