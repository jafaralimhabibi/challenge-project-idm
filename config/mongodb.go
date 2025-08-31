package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client

func InitMongo(ctx context.Context) (*mongo.Client, error) {
	if mongoClient != nil {
		return mongoClient, nil
	}

	c := GetConfig()
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		c.MongoDB.Username,
		c.MongoDB.Password,
		c.MongoDB.Host,
		c.MongoDB.Port,
	)
	clientOpts := options.Client().ApplyURI(uri)
	if c.MongoDB.Username != "" {
		cred := options.Credential{
			Username: c.MongoDB.Username,
			Password: c.MongoDB.Password,
		}
		clientOpts.SetAuth(cred)
	}

	ctx2, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx2, clientOpts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx2, nil); err != nil {
		return nil, err
	}

	mongoClient = client
	log.Println("Connected to MongoDB")

	return mongoClient, nil
}

func GetDatabase() *mongo.Database {
	if mongoClient == nil {
		panic("Unknown mongoDB client")
	}
	return mongoClient.Database(GetConfig().MongoDB.Database)
}
