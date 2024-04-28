package store

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	requestTimeout = 3 * time.Second
)

type Storage struct {
	Config StorageConfig
	client *mongo.Client
	cursor *mongo.Cursor
}

func NewStorageConn(config StorageConfig) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	// Index crawled data by page URL
	model := mongo.IndexModel{
		Keys: bson.M{
			"url": 1, // ascending order
		},
		Options: options.Index().SetUnique(true),
	}
	collection := client.Database(config.DatabaseName).Collection(config.CrawledDataCollection)
	_, err = collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		return nil, err
	}

	// Index webgraph by URL
	model = mongo.IndexModel{
		Keys: bson.M{
			"url": 1, // ascending order
		},
		Options: options.Index().SetUnique(true),
	}
	collection = client.Database(config.DatabaseName).Collection(config.WebgraphCollection)
	_, err = collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		return nil, err
	}

	// Index word index by word
	model = mongo.IndexModel{
		Keys: bson.M{
			"word": 1, // ascending order
		},
		Options: options.Index().SetUnique(true),
	}
	collection = client.Database(config.DatabaseName).Collection(config.WordIndexCollection)
	_, err = collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		return nil, err
	}

	return &Storage{
		client: client,
		Config: config,
	}, nil
}

// Destroy disconnects from the database
func (db *Storage) Destroy() {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	if err := db.client.Disconnect(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to disconnect from MongoDB")
	}
}

// InitIterator initialises an iterator for the crawled data collection.
func (db *Storage) InitIterator(collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	collection := db.client.Database(db.Config.DatabaseName).Collection(collectionName)

	var err error
	db.cursor, err = collection.Find(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("Failed to read from collection %s: %v", db.Config.DatabaseName, err)
	}
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
