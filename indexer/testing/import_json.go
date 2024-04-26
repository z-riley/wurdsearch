package testing

import (
	"fmt"
	"os/exec"
)

const (
	databaseName   = "turdsearch"
	collectionName = "crawled_data_test"
	f              = "/home/zac/repo/turdsearch/indexer/testing/turdsearch.crawled_data.json"
)

// ImportJsonData imports the contents of a file into a collection
func ImportJsonData(filepath, database, collection string) error {

	cmd := exec.Command("mongoimport", "--db", database, "--collection", collection, "--file", filepath, "--jsonArray")

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error running mongoimport: %v", err)
	}
	return nil
}
