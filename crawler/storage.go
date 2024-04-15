package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type storage struct {
	client *mongo.Client
}

func newMongoDBConn() (*storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	return &storage{
		client: client,
	}, nil
}

func (db *storage) connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal().Err(err)
		}
	}()
}

func (db *storage) insert() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := db.client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal().Err(err)
	}
}

func (db *storage) destroy() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.client.Disconnect(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed ot disconnect from MongoDB")
	}
}
