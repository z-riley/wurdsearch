package store

import "time"

const (
	DatabaseName = "wurdsearch"

	CrawledDataCollection     = "crawled_data"
	CrawledDataTestCollection = "crawled_data_test"

	WebgraphCollection     = "webgraph"
	WebgraphTestCollection = "webgraph_test"

	WordIndexCollection     = "word_index"
	WordIndexTestCollection = "word_index_test"

	requestTimeout = 3 * time.Second
	connectionPool = 800 // manually tuned for best results
)

// Config holds settings used in the store package.
type Config struct {
	DatabaseName          string
	CrawledDataCollection string
	WebgraphCollection    string
	WordIndexCollection   string
}
