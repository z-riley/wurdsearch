package main

// wordIndex contains which websites use a particular word, and how many times
// it appears on each page
type wordIndex struct {
	word       string          `bson:"word"`
	references map[string]uint `bson:"references"`
}

// digest processes page data
func digest() error {
	return nil
}
