package testing

import (
	"testing"

	"github.com/z-riley/wurdsearch/common/store"
)

func TestImportJson(t *testing.T) {
	const (
		databaseName   = store.DatabaseName
		collectionName = store.CrawledDataTestCollection
		f              = "/home/zac/repo/wurdsearch/indexer/testing/wurdsearch.crawled_data.json"
	)

	if err := ImportJson(f, databaseName, collectionName); err != nil {
		t.Error(err)
	}
}
