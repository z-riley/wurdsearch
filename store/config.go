package store

const (
	DatabaseName = "turdsearch"

	CrawledDataCollection     = "crawled_data"
	CrawledDataTestCollection = "crawled_data_test"

	WebgraphCollection     = "webgraph"
	WebgraphTestCollection = "webgraph_test"

	WordIndexCollection     = "word_index"
	WordIndexTestCollection = "word_index_test"
)

type StorageConfig struct {
	DatabaseName          string
	CrawledDataCollection string
	WebgraphCollection    string
	WordIndexCollection   string
}
