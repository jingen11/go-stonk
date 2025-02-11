package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Init(url string) (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	timeoutCtx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(timeoutCtx, opts)
	if err != nil {
		return nil, err
	}
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(timeoutCtx, bson.D{primitive.E{Key: "ping", Value: 1}}).Err(); err != nil {
		return nil, err
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client, nil
}

func Disconnect(client *mongo.Client) {
	fmt.Println("Disconnecting from MongoDB!")
	client.Disconnect(context.Background())
}

func InitPriceCollection(db *mongo.Database) (*mongo.Collection, error) {
	priceColl := db.Collection("price")
	_, err := priceColl.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "symbol", Value: 1}, {Key: "date", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return nil, err
	}

	return priceColl, nil
}

func InitSymbolCollection(db *mongo.Database) (*mongo.Collection, error) {
	symbolColl := db.Collection("symbol")

	_, err := symbolColl.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys:    bson.D{{Key: "symbol", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return nil, err
	}

	return symbolColl, nil
}
