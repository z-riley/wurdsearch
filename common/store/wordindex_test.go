package store

import (
	"fmt"
	"testing"
)

func TestGetWordIndex(t *testing.T) {
	db, err := NewStorageConn(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()

	word, err := db.GetWordIndex("testword")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", word)
}

func TestUpdateWordReference(t *testing.T) {
	db, err := NewStorageConn(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()

	if err := db.UpdateWordReference("testword", "example.com", 1, 1000); err != nil {
		t.Fatal(err)
	}
}
