package testing

import (
	"fmt"
	"os/exec"

	"github.com/z-riley/turdsearch/common/store"
)

// ImportJson imports the contents of a file into a collection
func ImportJson(filepath, database, collection string) error {

	// Make a new storage object to index collection by URL
	db, err := store.NewStorageConn(store.StorageConfig{
		DatabaseName:          database,
		CrawledDataCollection: store.CrawledDataTestCollection,
		WebgraphCollection:    store.WebgraphTestCollection,
	})
	defer db.Destroy()

	cmd := exec.Command("mongoimport", "--db", database, "--collection", collection, "--file", filepath, "--jsonArray")
	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error running mongoimport: %v", err)
	}
	return nil
}
