package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Connector handles connection operation to the database
// TODO: handle retry mechanism efficiently
type Connector struct {
	uri string
}

//Connect connects to the database instance
func (c *Connector) Connect(uri string, option *options.ClientOptions) (*mongo.Client, error) {

	c.uri = uri
	client, err := mongo.NewClient(option.ApplyURI(c.uri))
	if err != nil {
		log.Error("Error connecting database: ", err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error("Error pinging database: ", err)
		return nil, err
	}
	return client, err
}

//Retry retries connection to mongoDB
func (c *Connector) Retry() (*mongo.Client, error) {

	client, err := mongo.NewClient(options.Client().ApplyURI(c.uri))
	if err != nil {
		log.Error("Error connecting database: ", err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client.Connect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error("Error connecting database: ", err)
		return nil, err
	}
	return client, err
}
