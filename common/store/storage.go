package store

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	Config StorageConfig
	client *mongo.Client
	cursor *mongo.Cursor
}

func NewStorageConn(config StorageConfig) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	clientOptions.SetMaxPoolSize(connectionPool)
	client, err := mongo.Connect(ctx, clientOptions)
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

// Len returns the number of documents in a collection
func (db *Storage) Len(collectionName string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	collection := db.client.Database(db.Config.DatabaseName).Collection(collectionName)
	// Improve performance by using a hint to take advantage of the built-in index on
	// the _id field. See https://www.mongodb.com/docs/drivers/go/upcoming/fundamentals/crud/read-operations/count/
	opts := options.Count().SetHint("_id_")
	length, err := collection.CountDocuments(ctx, bson.D{}, opts)
	if err != nil {
		return 0, err
	}

	return length, nil
}

// MaxConnections returns the size of the connection pool
func MaxConnections() int {
	return connectionPool
}
