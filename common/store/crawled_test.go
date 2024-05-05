package store

import (
	"fmt"
	"testing"
	"time"
)

func TestSavePageData(t *testing.T) {
	db, err := NewStorageConn(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()

	for _, data := range []PageData{
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
		if err := db.SavePageData(data); err != nil {
			t.Fatal(err)
		}
	}

}

func TestFetchPageData(t *testing.T) {
	db, err := NewStorageConn(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()

	url := "https://example.com/"
	data, err := db.FetchPageData(url)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("t: %+v\n", data)
}

func TestPageIsRecentlyCrawled(t *testing.T) {
	db, err := NewStorageConn(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()

	url := "https://example.com/"
	window := 24 * time.Hour
	result, err := db.PageIsRecentlyCrawled(url, window)
	if err != nil {
		t.Fatal(t)
	}

	fmt.Printf("Page %s crawled in the last %v: %t\n", url, window, result)
}

func TestPageLastCrawled(t *testing.T) {
	db, err := NewStorageConn(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()

	url := "https://example.com/"
	timeLastCrawled, err := db.PageLastCrawled(url)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Page %s last crawled at %v", url, timeLastCrawled)
}

func TestNextPageData(t *testing.T) {
	db, err := NewStorageConn(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()

	if err := db.InitIterator(CrawledDataTestCollection); err != nil {
		t.Fatal(err)
	}

	isMoreData := true
	for isMoreData {
		var data PageData
		var err error
		data, isMoreData, err = db.NextPageData()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%v\n", data.Url)
	}

}
