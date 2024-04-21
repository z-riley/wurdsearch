package main

import (
	"fmt"
	"testing"
	"time"
)

func TestNewStorageConn(t *testing.T) {
	db, err := newStorageConn(getTestConfig())
	if err != nil {
		t.Error(err)
	}
	defer db.destroy()
}

func TestPing(t *testing.T) {
	db, err := newStorageConn(getTestConfig())
	if err != nil {
		t.Error(err)
	}
	defer db.destroy()

	if err := db.ping(); err != nil {
		t.Error(err)
	}
}

func TestEnterPageData(t *testing.T) {
	db, err := newStorageConn(getTestConfig())
	if err != nil {
		t.Error(err)
	}
	defer db.destroy()

	for _, data := range []pageData{
		{
			Url:          "https://example.com/",
			LastAccessed: time.Now(),
			Links:        []string{"a link", "another link"},
			Content:      "lots of content yes",
		},
		{
			Url:          "https://anotherexample.com/",
			LastAccessed: time.Now(),
			Links:        []string{"more links", "another link"},
			Content:      "all my yummy content",
		},
	} {
		if err := db.savePageData(data); err != nil {
			t.Error(err)
		}
	}

}

func TestFetchPageData(t *testing.T) {
	db, err := newStorageConn(getTestConfig())
	if err != nil {
		t.Error(err)
	}
	defer db.destroy()

	url := "https://example.com/"
	data, err := db.fetchPageData(url)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("t: %+v\n", data)
}

func TestPageIsRecentlyCrawled(t *testing.T) {
	db, err := newStorageConn(getTestConfig())
	if err != nil {
		t.Error(err)
	}
	defer db.destroy()

	url := "https://example.com/"
	window := 24 * time.Hour
	result, err := db.pageIsRecentlyCrawled(url, window)
	if err != nil {
		t.Error(t)
	}

	fmt.Printf("Page %s crawled in the last %v: %t\n", url, window, result)
}

func TestPageLastCrawled(t *testing.T) {
	db, err := newStorageConn(getTestConfig())
	if err != nil {
		t.Error(err)
	}
	defer db.destroy()

	url := "https://example.com/"
	timeLastCrawled, err := db.pageLastCrawled(url)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("Page %s last crawled at %v", url, timeLastCrawled)
}

func getTestConfig() storageConfig {
	// TODO: use a seperate DB for testing
	return storageConfig{
		databaseName:          "turdsearch",
		crawledDataCollection: "crawled_data",
		indexedDataCollection: "indexed_data",
	}
}
