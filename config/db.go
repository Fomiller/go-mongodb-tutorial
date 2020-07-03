package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CLIENT *mongo.Client

func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// connect to MongoDB
	var err error
	CLIENT, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = CLIENT.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	// Log message if successful
	fmt.Println("Connected to MongoDB!")
}
