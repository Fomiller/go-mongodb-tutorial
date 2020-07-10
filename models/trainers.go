package models

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/fomiller/go-mongodb-tutorial/config"
)

var (
	collection *mongo.Collection
)

type Trainer struct {
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Age  int    `json:"age,omitempty" bson:"age,omitempty"`
	City string `json:"city,omitempty" bson:"city,omitempty"`
}

func init() {
	// create trainer collection
	collection = config.CLIENT.Database("go-mongo-tut").Collection("trainers")
}

func AllTrainers() []*Trainer {
	findOptions := options.Find()
	findOptions.SetLimit(10)

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
	return results
}

func OneTrainer() {

}

func CreateTrainer(req *http.Request) {
	var err error
	var newTrainer Trainer
	newTrainer.Name = req.FormValue("name")
	newTrainer.Age, err = strconv.Atoi(req.FormValue("age"))
	if err != nil {
		log.Panic(err)
	}
	newTrainer.City = req.FormValue("city")

	insertResult, err := collection.InsertOne(context.TODO(), newTrainer)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println("inserted a single document: ", insertResult)
}
