package store

const (
	DatabaseName = "turdsearch"

	CrawledDataCollection     = "crawled_data"
	CrawledDataTestCollection = "crawled_data_test"

	WebgraphCollection     = "webgraph"
	WebgraphTestCollection = "webgraph_test"
)

type StorageConfig struct {
	DatabaseName, CrawledDataCollection, WebgraphCollection string
}
