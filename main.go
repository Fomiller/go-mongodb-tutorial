package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/fomiller/go-mongodb-tutorial/API"
	"github.com/fomiller/go-mongodb-tutorial/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string
	Age  int
	City string
}

func main() {

	// create collection trainer
	collection := config.CLIENT.Database("go-mongo-tut").Collection("trainers")

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

	// create filter for document with the name of "Ash"
	filter := bson.D{{Key: "name", Value: "Ash"}}

	// Update the document to increment($inc) the age of the document by "2"
	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "age", Value: 1},
		}},
	}

	// update collection
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Find
	// create variable to be pointed at when returning result from query
	var result Trainer

	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)

	// find multiple
	// pass these options to the Find Method
	findOptions := options.Find()
	findOptions.SetLimit(2)

	// Create Slice to store decoded documents in
	var results []*Trainer

	// passing bson.D{{}} as the filer matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem Trainer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// close the cursor once finished
	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of poointers): %+v\n", results)

	// delete documents
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v Document Deleted", deleteResult.DeletedCount)

	// Route Handlers
	http.HandleFunc("/", API.IndexHandler)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.ListenAndServe(":8080", nil)
}
