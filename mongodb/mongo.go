package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client    *mongo.Client
	mongoOnce sync.Once
)

type MongoConfig struct {
	URI string `json:"uri"`
}

type MongoClient struct{}

func NewMongoClient(conf MongoConfig) (*MongoClient, error) {
	var err error
	mongoOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(conf.URI))
		if err != nil {
			return
		}
	})
	return &MongoClient{}, err
}

func (c *MongoClient) GetID(dbName, collection string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result struct {
		Name   string `json:"name" bson:"name"`
		NextID int64  `json:"next_id" bson:"next_id"`
	}

	err := client.Database(dbName).Collection(ID).FindOneAndUpdate(
		ctx,
		bson.M{"name": collection},
		bson.M{"$inc": bson.M{"next_id": 1}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		_, err = client.Database(dbName).Collection(ID).InsertOne(ctx, bson.M{"name": collection, "next_id": 1})
		if err != nil {
			return -1, err
		}
		return 1, err
	} else if err != nil {
		return -1, err
	}
	return result.NextID, err
}
