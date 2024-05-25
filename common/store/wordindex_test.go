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

func TestUpdateWordReferences(t *testing.T) {
	db, err := NewStorageConn(getTestConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Destroy()

	wordCounts := map[string]uint{
		"please":  2,
		"disable": 4,
		"nuts":    8,
	}
	if err := db.UpdateWordReferences("example.com", wordCounts, 1000); err != nil {
		t.Fatal(err)
	}
}
