package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// WordEntry contains which websites use a particular word, and how many times
// it appears on each page.
type WordEntry struct {
	// Word is the word in the entry
	Word string `bson:"word"`
	// References maps each URL to the number of occurrances and total length of the document
	References map[string]Reference `bson:"references"`
}

// Encode changes "."s to "`"s because "." confuses Mongo
func (w *WordEntry) Encode() WordEntry {
	encodedRefs := make(map[string]Reference)
	for url, reference := range w.References {
		encodedURL := strings.ReplaceAll(url, ".", "`")
		encodedRefs[encodedURL] = reference
	}
	return WordEntry{
		Word:       w.Word,
		References: encodedRefs,
	}
}

// Decode changes "`"s to "."s because "." confuses poor Mongo.
func (w *WordEntry) Decode() WordEntry {
	decodedRefs := make(map[string]Reference)
	for url, reference := range w.References {
		encodedURL := strings.ReplaceAll(url, "`", ".")
		decodedRefs[encodedURL] = reference
	}
	return WordEntry{
		Word:       w.Word,
		References: decodedRefs,
	}
}

type Reference struct {
	Count  uint `bson:"count"`
	Length uint `bson:"length"`
}

// Getword retrieves a word index document.
func (db *Storage) GetWordIndex(word string) (WordEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	collection := db.client.Database(db.Config.DatabaseName).Collection(db.Config.WordIndexCollection)
	var retrievedWordEntry WordEntry

	err := collection.FindOne(ctx, bson.M{"word": word}).Decode(&retrievedWordEntry)
	if err != nil {
		return WordEntry{}, err
	}
	return retrievedWordEntry.Decode(), nil
}

// SaveWord inserts or overwrites a word index doc.
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
		// log.Trace().Msgf("Inserted new DB document for word: %s", word.Word)
	} else if result.ModifiedCount == 1 {
		// log.Trace().Msgf("Updated existing DB document for word: %s", word.Word)
	}

	return nil
}

// UpdateWordReferences inserts or overwrites a word index document for a given URL.
func (db *Storage) UpdateWordReferences(URL string, wordCounts map[string]uint, totalWords uint) error {
	// Prepare models for bulk write
	var models []mongo.WriteModel
	for word, count := range wordCounts {
		model := mongo.NewUpdateOneModel().
			SetFilter(bson.D{{Key: "word", Value: word}}).
			SetUpdate(bson.D{{
				Key: "$set",
				Value: bson.D{{
					Key:   fmt.Sprintf("references.%s", strings.ReplaceAll(URL, ".", "`")),
					Value: Reference{count, totalWords},
				}},
			}}).
			SetUpsert(true)
		models = append(models, model)
	}

	// Update documents
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	collection := db.client.Database(db.Config.DatabaseName).Collection(db.Config.WordIndexCollection)
	opts := options.BulkWrite().SetOrdered(true)
	results, err := collection.BulkWrite(ctx, models, opts)
	if err != nil {
		return fmt.Errorf("Failed to bulk write: %v", err)
	}

	log.Debug().Msgf("Matched %d documents, modified %d documents, upserted %d documents", results.MatchedCount, results.ModifiedCount, results.UpsertedCount)

	return nil
}

// NextPageData gets the next word entry document. InitIterator must be called first.
// Returns true if there is more data to iterate over
func (db *Storage) NextWordEntry() (WordEntry, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	if !db.cursor.Next(ctx) {
		if err := db.cursor.Err(); err != nil {
			return WordEntry{}, false, fmt.Errorf("Cursor error: %v", err)
		}
		return WordEntry{}, false, nil
	}

	var result WordEntry
	if err := db.cursor.Decode(&result); err != nil {
		return result, true, fmt.Errorf("Failed to decode word entry: %v", err)
	}

	return result, true, nil
}
