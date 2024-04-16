package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	databaseName          = "turdsearch"
	crawledDataCollection = "crawled_data"
	indexedDataCollection = "indexed_data"

	requestTimeout = 3 * time.Second
)

type pageData struct {
	Url          string    `bson:"url"`
	LastAccessed time.Time `bson:"lastAccessed"`
	Links        []string  `bson:"links"`
	Content      string    `bson:"content"`
}

type storage struct {
	client *mongo.Client
}

func newStorageConn() (*storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	// Index by page URL
	model := mongo.IndexModel{
		Keys: bson.M{
			"url": 1, // ascending order
		},
		Options: options.Index().SetUnique(true),
	}

	// Create the index
	collection := client.Database(databaseName).Collection(crawledDataCollection)
	_, err = collection.Indexes().CreateOne(ctx, model)
	if err != nil {
		return nil, err
	}

	return &storage{
		client: client,
	}, nil
}

func (db *storage) ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	err := db.client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}
	return nil
}

// savePageData inserts or overwrites a page data document
func (db *storage) savePageData(data pageData) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	filter := bson.M{"url": data.Url}        // match document by URL (URL is already defined as a unique index)
	update := bson.M{"$set": data}           // set new data for the document
	opts := options.Update().SetUpsert(true) // create a new document if one doesn't exist

	collection := db.client.Database(databaseName).Collection(crawledDataCollection)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 && result.UpsertedCount == 1 {
		log.Debug().Msgf("Inserted new document for URL: %s", data.Url)
	} else if result.ModifiedCount == 1 {
		log.Debug().Msgf("Updated existing document for URL: %s", data.Url)
	}

	return nil
}

// fetchPageData retrieves page data for a specified URL
func (db *storage) fetchPageData(url string) (pageData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	collection := db.client.Database(databaseName).Collection(crawledDataCollection)
	var retrievedPageData pageData

	err := collection.FindOne(ctx, bson.M{"url": url}).Decode(&retrievedPageData)
	if err != nil {
		return pageData{}, err
	}
	log.Info().Msgf("Retrieved event date: %+v\n", retrievedPageData)
	return retrievedPageData, nil
}

func (db *storage) destroy() {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	if err := db.client.Disconnect(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to disconnect from MongoDB")
	}
}
