package main

import "testing"

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
