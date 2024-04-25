package main

import (
	"fmt"
	"os/exec"
)

const (
	databaseName   = "turdsearch"
	collectionName = "crawled_data_test"
)

func importTestDataJson(filepath ...string) error {
	if len(filepath) == 0 {
		filepath = []string{"/home/zac/repo/turdsearch/indexer/turdsearch.crawled_data.json"}
	}
	cmd := exec.Command("mongoimport", "--db", databaseName, "--collection", collectionName, "--file", filepath[0], "--jsonArray")

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error running mongoimport: %v", err)
	}
	return nil
}
