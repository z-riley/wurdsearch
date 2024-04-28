package store

import (
	"fmt"
	"testing"
)

func TestNewStorageConn(t *testing.T) {
	db, err := NewStorageConn(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()
}

func TestLen(t *testing.T) {
	db, err := NewStorageConn(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()

	length, err := db.Len(CrawledDataTestCollection)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(length)
}

func getTestConfig() StorageConfig {
	return StorageConfig{
		DatabaseName:          DatabaseName,
		CrawledDataCollection: CrawledDataTestCollection,
		WebgraphCollection:    WebgraphTestCollection,
		WordIndexCollection:   WordIndexTestCollection,
	}
}
