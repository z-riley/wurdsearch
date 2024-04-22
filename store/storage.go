package store

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	requestTimeout = 3 * time.Second
)

type PageData struct {
	Url          string    `bson:"url"`
	LastAccessed time.Time `bson:"lastAccessed"`
	Links        []string  `bson:"links"`
	Content      string    `bson:"content"`
}

type Storage struct {
	client *mongo.Client
	Config StorageConfig
}

type StorageConfig struct {
	DatabaseName, CrawledDataCollection, IndexedDataCollection string
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

	return &Storage{
		client: client,
		Config: config,
	}, nil
}

// SavePageData inserts or overwrites a page data document
func (db *Storage) SavePageData(data PageData) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	filter := bson.M{"url": data.Url}        // match document by URL (URL is already defined as a unique index)
	update := bson.M{"$set": data}           // set new data for the document
	opts := options.Update().SetUpsert(true) // create a new document if one doesn't exist

	collection := db.client.Database(db.Config.DatabaseName).Collection(db.Config.CrawledDataCollection)
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 && result.UpsertedCount == 1 {
		log.Debug().Msgf("Inserted new DB document for URL: %s", data.Url)
	} else if result.ModifiedCount == 1 {
		log.Debug().Msgf("Updated existing DB document for URL: %s", data.Url)
	}

	return nil
}

// FetchPageData retrieves page data for a specified URL
func (db *Storage) FetchPageData(url string) (PageData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	collection := db.client.Database(db.Config.DatabaseName).Collection(db.Config.CrawledDataCollection)
	var retrievedPageData PageData

	err := collection.FindOne(ctx, bson.M{"url": url}).Decode(&retrievedPageData)
	if err != nil {
		return PageData{}, err
	}
	log.Info().Msgf("Retrieved event date: %+v\n", retrievedPageData)
	return retrievedPageData, nil
}

// PageIsRecentlyCrawled checks if a page was crawled in the last specified time frame (window)
func (db *Storage) PageIsRecentlyCrawled(url string, window time.Duration) (bool, error) {
	lastCrawled, err := db.PageLastCrawled(url)
	if err != nil {
		return false, err
	}

	if time.Since(lastCrawled) > window {
		return false, nil
	} else {
		return true, nil
	}
}

// PageLastCrawled returns when a page was last crawled. Error if page doesn't exist
func (db *Storage) PageLastCrawled(url string) (time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	collection := db.client.Database(db.Config.DatabaseName).Collection(db.Config.CrawledDataCollection)

	var result bson.M
	err := collection.FindOne(ctx, bson.M{"url": url}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return time.Unix(0, 0), nil
	} else if err != nil {
		return time.Unix(0, 0), err
	}

	field := "lastAccessed"
	if dt, ok := result[field].(primitive.DateTime); ok {
		return dt.Time(), nil
	} else {
		return time.Unix(0, 0), fmt.Errorf("Failed to convert field %s into valid time", field)
	}
}

// Destroy disconnects from the database
func (db *Storage) Destroy() {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	if err := db.client.Disconnect(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to disconnect from MongoDB")
	}
}
