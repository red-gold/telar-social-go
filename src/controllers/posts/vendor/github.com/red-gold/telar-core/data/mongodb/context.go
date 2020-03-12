package mongodb

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/red-gold/telar-core/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Struct used for maintaining HTTP Request Context
type MongoClient struct {
	Context context.Context
}

var client *mongo.Client

// Create a new mogodb client object
func NewMongoClient(ctx context.Context) (MongoDatabase, error) {
	mongoClient := &MongoClient{Context: ctx}

	mongoPwd, readFileErr := ioutil.ReadFile(*config.AppConfig.MongoPwdPath)
	if readFileErr != nil {
		log.Fatalf("unable to read path: %s, error: %s", *config.AppConfig.MongoPwdPath, readFileErr.Error())
	}

	url := fmt.Sprintf(*config.AppConfig.MongoDBHost, mongoPwd)

	fmt.Printf("Mongo URL: %s", url)
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(url))

	return mongoClient, err
}

// GetContext get mongodb context
func (c *MongoClient) GetContext() (context.Context, error) {
	if client == nil {
		return nil, fmt.Errorf("MongoDB client is nil, be sure you have invoked NewClient() function already!")
	}
	return c.Context, nil
}

// Close Client
func (c *MongoClient) Close() error {
	return client.Disconnect(c.Context)
}

// GetCollection gets database collection
func (c *MongoClient) GetCollection(collectionName string) (*mongo.Collection, error) {
	if client == nil {
		return nil, fmt.Errorf("MongoDB client is nil, be sure you have invoked NewClient() function already!")
	}
	collection := client.Database(*config.AppConfig.Database).Collection(collectionName)

	return collection, nil
}

// Create database session
func (c *MongoClient) GetDb() (*mongo.Database, error) {
	if client == nil {
		return nil, fmt.Errorf("MongoDB client is nil, be sure you have invoked NewClient() function already!")
	}
	db := client.Database(*config.AppConfig.Database)

	return db, nil
}

// Ping database
func (c *MongoClient) Ping() error {
	if client == nil {
		_ = client.Ping(c.Context, readpref.Primary())
	}

	return nil
}
