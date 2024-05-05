package store

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Node contains data about which other sites link to and from itself.
// This is used to help calculate the "importance" of the site
type Node struct {
	Url       string   `bson:"url"`
	LinksTo   []string `bson:"linksTo"`
	LinksFrom []string `bson:"linksFrom"`
}

// SaveNode inserts or overwrites a webgraph node
func (db *Storage) SaveNode(node Node) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	filter := bson.M{"url": node.Url}        // match document by URL (URL is already defined as a unique index)
	update := bson.M{"$set": node}           // set new data for the document
	opts := options.Update().SetUpsert(true) // create a new document if one doesn't exist

	collection := db.client.Database(db.Config.DatabaseName).Collection(db.Config.WebgraphCollection)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 && result.UpsertedCount == 1 {
		log.Debug().Msgf("Inserted new DB document for URL: %s", node.Url)
	} else if result.ModifiedCount == 1 {
		log.Debug().Msgf("Updated existing DB document for URL: %s", node.Url)
	}

	return nil
}
