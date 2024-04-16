package main

import (
	"fmt"
	"testing"
	"time"
)

func TestNewStorageConn(t *testing.T) {
	db, err := newStorageConn()
	if err != nil {
		t.Error(err)
	}
	defer db.destroy()
}

func TestPing(t *testing.T) {
	db, err := newStorageConn()
	if err != nil {
		t.Error(err)
	}
	defer db.destroy()

	if err := db.ping(); err != nil {
		t.Error(err)
	}
}

func TestEnterPageData(t *testing.T) {
	db, err := newStorageConn()
	if err != nil {
		t.Error(err)
	}
	defer db.destroy()

	// data := pageData{
	// 	Url:          "https://example.com/",
	// 	LastAccessed: time.Now(),
	// 	Links:        []string{"a link", "another link"},
	// 	Content:      "lots of content yes",
	// }

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
		if err := db.enterPageData(data); err != nil {
			t.Error(err)
		}
	}

}

func TestFetchPageData(t *testing.T) {
	db, err := newStorageConn()
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
