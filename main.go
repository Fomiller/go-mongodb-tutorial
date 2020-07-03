package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/fomiller/go-mongodb-tutorial/API"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	// Log message if successful
	fmt.Println("Connected to MongoDB!")

	// create collection trainer
	collection := client.Database("go-mongo-tut").Collection("trainers")

	// create trainers to add into database
	ash := Trainer{"Ash", 10, "Pallet Town"}
	misty := Trainer{"Misty", 10, "Cerulean City"}
	brock := Trainer{"Brock", 15, "Pewter City"}

	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}
	// successful insert
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// create slice of trainers to insert at once
	trainers := []interface{}{misty, brock}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted multiple documents: ", insertManyResult.InsertedIDs)

	// Route Handlers
	http.HandleFunc("/", API.IndexHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}
