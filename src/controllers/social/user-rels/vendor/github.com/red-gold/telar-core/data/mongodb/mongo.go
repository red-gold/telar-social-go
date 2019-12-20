package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDatabase interface {
	GetCollection(collectionName string) (*mongo.Collection, error)
	GetDb() (*mongo.Database, error)
	Close() error
	GetContext() (context.Context, error)
	Ping() error
}
