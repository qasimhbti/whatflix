package main

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const dbName = "whatflix"

func newMongoDBClientGetter(cs string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(cs))
	if err != nil {
		return nil, errors.Wrap(err, "new client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "connect")
	}
	log.Println("mongo DB connected")

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, errors.WithMessage(err, "ping")
	}
	log.Println("mongo DB : ping successfully")
	return client, nil
}

func getmgoDB(client *mongo.Client) *mongo.Database {
	return client.Database(dbName)
}
