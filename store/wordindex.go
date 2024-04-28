package store

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// WordEntry contains which websites use a particular word, and how many times
// it appears on each page
type WordEntry struct {
	Word       string          `bson:"word"`
	References map[string]uint `bson:"references"`
}

// Getword retrieves a word index document
func (db *Storage) GetWord(word string) (WordEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	collection := db.client.Database(db.Config.DatabaseName).Collection(db.Config.WordIndexCollection)
	var retrievedWordEntry WordEntry

	err := collection.FindOne(ctx, bson.M{"word": word}).Decode(&retrievedWordEntry)
	if err != nil {
		return WordEntry{}, err
	}
	return retrievedWordEntry, nil
}

// SaveWord inserts or overwrites a word index doc
func (db *Storage) SaveWord(word WordEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	filter := bson.M{"word": word.Word}      // match document by word (word is already defined as a unique index)
	update := bson.M{"$set": word}           // set new data for the document
	opts := options.Update().SetUpsert(true) // create a new document if one doesn't exist

	collection := db.client.Database(db.Config.DatabaseName).Collection(db.Config.WordIndexCollection)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 && result.UpsertedCount == 1 {
		log.Debug().Msgf("Inserted new DB document for word: %s", word.Word)
	} else if result.ModifiedCount == 1 {
		log.Debug().Msgf("Updated existing DB document for word: %s", word.Word)
	}

	return nil
}

// UpdateWordReference inserts or overwrites a word index document. Only one "url: count" reference may be added at a time
func (db *Storage) UpdateWordReference(word string, url string, count uint) error {
	// Manually read and overwrite document because Mongo can't handle "." characters in keys
	wordEntry, err := db.GetWord(word)
	if errors.Is(err, mongo.ErrNoDocuments) {
		// If no doc for that word exists, make one
		wordEntry = WordEntry{
			Word: word,
			References: map[string]uint{
				url: count,
			},
		}
	} else if err != nil {
		return err
	}
	wordEntry.References[url] = count
	if err := db.SaveWord(wordEntry); err != nil {
		return err
	}

	return nil
}
