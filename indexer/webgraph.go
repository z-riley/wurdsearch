package main

// node contains data about which other sites link to and from itself.
// This is used to help calculate the "importance" of the site
type node struct {
	url       string   `bson:"url"`
	linksTo   []string `bson:"linksTo"`
	linksFrom []string `bson:"linksFrom"`
}

type Webgrapher struct{}

func NewWebgrapher() *Webgrapher {
	return &Webgrapher{}
}

// GenerateWebgraph generates a webgraph from the crawled data in the database
func (w *Webgrapher) GenerateWebgraph() error {
	return nil
}
