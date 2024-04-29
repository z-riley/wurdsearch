package search

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// vector holds an n-dimensional vector with with value for each dimension
type vector struct {
	label string
	val   map[string]float64
}

// mod calculates the modulus (magnitude) of a vector
func (v *vector) mod() float64 {
	sum := 0.0
	for _, val := range v.val {
		sum += (val * val)
	}
	return math.Sqrt(sum)
}

// theta returns the angle (in radians) between two vectors
func theta(a vector, b vector) float64 {
	return math.Acos(dotProduct(a, b) / (a.mod() * b.mod()))
}

// dotProduct calculates the dot product of two vectors. It is assumed that vector.vals contain the same keys
func dotProduct(a vector, b vector) float64 {
	sum := 0.0
	for word := range a.val {
		sum += a.val[word] * b.val[word]
	}
	return sum
}

func (s *Searcher) TFIDF(searchTerm string) error {

	// Calculate search vectors
	searchTerm = strings.ToLower(searchTerm)
	searchTerm = strings.TrimSpace(searchTerm)
	words := strings.Split(searchTerm, " ")
	urls, err := s.getEveryRelevantDoc(words)
	if err != nil {
		return err
	}
	var vectors []vector
	for _, url := range urls {
		v, err := s.generateVector(url, words)
		if err != nil {
			return err
		}
		vectors = append(vectors, v)
	}
	fmt.Println(vectors)

	// 2. Get search term vector
	searchVec, err := s.searchTermVector(words)
	if err != nil {
		return err
	}
	fmt.Println(searchVec)

	// 3. Compare each vectors to search term vector

	return nil
}

// generateVector calcuates a search vector for a given search term for a page
func (s *Searcher) generateVector(url string, searchWords []string) (vector, error) {
	v := vector{
		label: url,
		val:   make(map[string]float64),
	}

	for _, word := range searchWords {

		TFs, err := s.termFrequencies(word)
		if err != nil {
			return v, err
		}
		TF := TFs[url]

		// Calculate IDF so rarer words are more important
		IDF, err := s.inverseDocumentFrequency(word)
		if err != nil {
			return v, err
		}

		v.val[word] = TF * IDF
	}

	return v, nil
}

// searchTermVector gets the TF-IDF vector the search term
func (s *Searcher) searchTermVector(words []string) (vector, error) {
	result := vector{
		label: "searchTerm",
		val:   map[string]float64{},
	}

	// Count occurances of each word
	wordCounts := make(map[string]uint)
	for _, word := range words {
		wordCounts[word] += 1
	}

	// Get TF-IDF of each word in search term
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

// getEveryRelevantDoc retrieves the URL of every document that contains any of the words in the search term
func (s *Searcher) getEveryRelevantDoc(words []string) ([]string, error) {
	var urls []string
	for _, word := range words {
		wordIndexes, err := s.db.GetWordIndex(word)
		if err != nil {
			return urls, err
		}
		for url := range wordIndexes.References {
			urls = append(urls, url)
		}
	}
	return urls, nil
}

// pageTFs holds TF scores for URLs
type pageTFs map[string]float64

// termFrequencies returns the TFs of a specified word in every document.
// The TF is the number of times a word appears in a document divided by the total number
// of words in the document
func (s *Searcher) termFrequencies(word string) (pageTFs, error) {
	TFs := make(pageTFs)

	// Get all links which contain that word
	doc, err := s.db.GetWordIndex(word)
	if err != nil {
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
