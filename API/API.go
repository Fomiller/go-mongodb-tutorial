package API

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fomiller/go-mongodb-tutorial/config"
	"github.com/fomiller/go-mongodb-tutorial/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// create collection for trainers
	collection *mongo.Collection
	// create trainers to add into database
	ash   = models.Trainer{"Ash", 10, "Pallet Town"}
	misty = models.Trainer{"Misty", 10, "Cerulean City"}
	brock = models.Trainer{"Brock", 15, "Pewter City"}
	// create filter for document with the name of "Ash"
	filter = bson.D{{Key: "name", Value: "Ash"}}
)

func init() {
	// create collection trainer
	collection = config.CLIENT.Database("go-mongo-tut").Collection("trainers")
}

func IndexHandler(res http.ResponseWriter, req *http.Request) {
	config.TPL.ExecuteTemplate(res, "index.html", nil)

}

func CreateHandler(res http.ResponseWriter, req *http.Request) {
	insertResult, err := collection.InsertOne(context.TODO(), ash)
	if err != nil {
		log.Fatal(err)
	}

	// successful insert
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// respond with json
	json.NewEncoder(res).Encode(insertResult)

}

func CreateManyHandler(res http.ResponseWriter, req *http.Request) {
	// create slice of trainers to insert at once
	trainers := []interface{}{misty, brock}

	insertManyResult, err := collection.InsertMany(context.TODO(), trainers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("inserted multiple documents: ", insertManyResult.InsertedIDs)
	json.NewEncoder(res).Encode(insertManyResult)

}

func UpdateHandler(res http.ResponseWriter, req *http.Request) {

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
	json.NewEncoder(res).Encode(updateResult)

}

func FindHandler(res http.ResponseWriter, req *http.Request) {
	// Find
	// create variable to be pointed at when returning result from query
	var result models.Trainer

	if err := collection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)
	json.NewEncoder(res).Encode(result)
}

func FindManyHandler(res http.ResponseWriter, req *http.Request) {
	// find multiple
	// pass these options to the Find Method
	findOptions := options.Find()
	findOptions.SetLimit(2)

	// Create Slice to store decoded documents in
	var results []*models.Trainer

	// passing bson.D{{}} as the filer matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.Trainer
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
	json.NewEncoder(res).Encode(results)

}

func DeleteHandler(res http.ResponseWriter, req *http.Request) {
	// delete documents
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v Document Deleted", deleteResult.DeletedCount)

}
