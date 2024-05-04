package lemmatiser

import (
	"fmt"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
)

type Lemmatiser struct {
	lemmatiser *golem.Lemmatizer
}

// NewLemmatiser creates a new lemmatisor object
func NewLemmatiser() (*Lemmatiser, error) {
	lemmatiser, err := golem.New(en.New())
	if err != nil {
		return nil, fmt.Errorf("Failed to create new lemmatiser: %v", err)
	}
	return &Lemmatiser{
		lemmatiser: lemmatiser,
	}, nil
}

// Lemmatise returns a word in it's base form. E.g., "running" -> "run".
// Based on https://raw.githubusercontent.com/michmech/lemmatization-lists/master/lemmatization-en.txt
func (l *Lemmatiser) Lemmatise(word string) string {
	return l.lemmatiser.Lemma(word)
}
