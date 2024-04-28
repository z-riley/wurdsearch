package store

import (
	"testing"
)

func TestNewStorageConn(t *testing.T) {
	db, err := NewStorageConn(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()
}

func getTestConfig() StorageConfig {
	return StorageConfig{
		DatabaseName:          DatabaseName,
		CrawledDataCollection: CrawledDataTestCollection,
		WebgraphCollection:    WebgraphTestCollection,
		WordIndexCollection:   WordIndexTestCollection,
	}
}
