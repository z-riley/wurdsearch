package store

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// WordEntry contains which websites use a particular word, and how many times
// it appears on each page
type WordEntry struct {
	Word       string          `bson:"word"`
	References map[string]uint `bson:"references"`
}

// SaveWord inserts or overwrites a word document
func (db *Storage) SaveWord(wordEntry WordEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	filter := bson.M{"word": wordEntry.Word} // match document by word (word is already defined as a unique index)
	update := bson.M{"$set": wordEntry}      // set new data for the document
	opts := options.Update().SetUpsert(true) // create a new document if one doesn't exist

	collection := db.client.Database(db.Config.DatabaseName).Collection(db.Config.WordIndexCollection)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 && result.UpsertedCount == 1 {
		log.Debug().Msgf("Inserted new DB document for word: %s", wordEntry.Word)
	} else if result.ModifiedCount == 1 {
		log.Debug().Msgf("Updated existing DB document for word: %s", wordEntry.Word)
	}

	return nil
}
