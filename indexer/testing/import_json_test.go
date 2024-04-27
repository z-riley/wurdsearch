package testing

import "testing"

func TestImportJson(t *testing.T) {
	const (
		databaseName   = "turdsearch"
		collectionName = "crawled_data_test"
		f              = "/home/zac/repo/turdsearch/indexer/testing/turdsearch.crawled_data.json"
	)

	if err := ImportJson(f, databaseName, collectionName); err != nil {
		t.Error(err)
	}
}
