package main

import (
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

	data := pageData{
		Url:          "https://example.com/",
		LastAccessed: time.Now(),
		Links:        []string{"a link", "another link"},
		Content:      "lots of content yes",
	}

	if err := db.enterPageData(data); err != nil {
		t.Error(err)
	}
}

func TestFetchData(t *testing.T) {
	db, err := newStorageConn()
	if err != nil {
		t.Error(err)
	}
	defer db.destroy()

	if err := db.fetchData(); err != nil {
		t.Error(err)
	}
}
