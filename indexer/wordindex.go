package main

// WordEntry contains which websites use a particular word, and how many times
// it appears on each page
type WordEntry struct {
	word       string          `bson:"word"`
	references map[string]uint `bson:"references"`
}

type WordIndexer struct{}

func NewWordIndexer() *WordIndexer {
	return &WordIndexer{}
}
