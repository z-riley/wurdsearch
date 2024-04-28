package search

import "fmt"

func (s *Searcher) TFIDF(term string) (float64, error) {
	return 0.0, nil
}

// wordFrequency returns the number of times a word appears in a document
// as a proportion of the total number of words in the document
func (s *Searcher) wordFrequency(word string) (float64, error) {

	// 1. Get all links which contain that word
	data, err := s.db.GetWord(word)
	if err != nil {
		return 0.0, err
	}

	fmt.Println(data.References)

	return 0.1, nil
}

func (s *Searcher) inverseDocumentFrequency() {}
