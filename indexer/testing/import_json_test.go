package testing

import (
	"testing"

	"github.com/z-riley/turdsearch/common/store"
)

func TestImportJson(t *testing.T) {
	const (
		databaseName   = store.DatabaseName
		collectionName = store.CrawledDataTestCollection
		f              = "/home/zac/repo/turdsearch/indexer/testing/turdsearch.crawled_data.json"
	)

	if err := ImportJson(f, databaseName, collectionName); err != nil {
		t.Error(err)
	}
}
