package store

import "time"

const (
	DatabaseName = "turdsearch"

	CrawledDataCollection     = "crawled_data"
	CrawledDataTestCollection = "crawled_data_test"

	WebgraphCollection     = "webgraph"
	WebgraphTestCollection = "webgraph_test"

	WordIndexCollection     = "word_index"
	WordIndexTestCollection = "word_index_test"

	requestTimeout = 3 * time.Second
	connectionPool = 800 // manually tuned for best results
)

type StorageConfig struct {
	DatabaseName          string
	CrawledDataCollection string
	WebgraphCollection    string
	WordIndexCollection   string
}