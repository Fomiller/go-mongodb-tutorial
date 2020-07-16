package API

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fomiller/go-mongodb-tutorial/config"
	"github.com/fomiller/go-mongodb-tutorial/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

type Test struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
	City string `json:"city,omitempty"`
}

var jsonTest = `{"name":"Forrest", "age": 26, "city": "Franklin"`

func init() {
	// create collection trainer
	collection = config.CLIENT.Database("go-mongo-tut").Collection("trainers")
}

func CreateHandler(res http.ResponseWriter, req *http.Request) {
	// init trainer variable for decoding
	var err error
	var newTrainer models.Trainer
	newTrainer.Name = req.FormValue("name")
	newTrainer.Age, _ = strconv.Atoi(req.FormValue("age"))
	newTrainer.City = req.FormValue("city")

	// insert into database
	insertResult, err := collection.InsertOne(context.TODO(), newTrainer)
	if err != nil {
		log.Fatal(err)
	}

	// // successful insert
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	config.TPL.ExecuteTemplate(res, "updated.gohtml", newTrainer)
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
	// var result models.Trainer

	// if err := collection.FindOne(context.TODO(), filter).Decode(&result); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Printf("Found a single document: %+v\n", result)
	// json.NewEncoder(res).Encode(result)
	result := models.OneTrainer()
	// fmt.Printf("Found FROM API: %+v\n", result)
	json.NewEncoder(res).Encode(result)
}

func FindManyHandler(res http.ResponseWriter, req *http.Request) {
	// find multiple
	// pass these options to the Find Method
	// findOptions := options.Find()
	// findOptions.SetLimit(10)

	// // Create Slice to store decoded documents in
	// var results []*models.Trainer

	// // passing bson.D{{}} as the filer matches all documents in the collection
	// cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Finding multiple documents returns a cursor
	// // Iterating through the cursor allows us to decode documents one at a time
	// for cur.Next(context.TODO()) {
	// 	// create a value into which the single document can be decoded
	// 	var elem models.Trainer
	// 	err := cur.Decode(&elem)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	results = append(results, &elem)
	// }

	// if err := cur.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	// // close the cursor once finished
	// cur.Close(context.TODO())

	// fmt.Printf("Found multiple documents (array of poointers): %+v\n", results)
	allTrainers := models.AllTrainers()
	config.TPL.ExecuteTemplate(res, "trainers.gohtml", allTrainers)
	// json.NewEncoder(res).Encode(results)

}

func DeleteHandler(res http.ResponseWriter, req *http.Request) {
	// delete documents
	deleteResult, err := collection.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v Document Deleted", deleteResult.DeletedCount)

}

// This is all code that i was using to read and test parsing the request body in learning how to get data from form submission
// I am keeping this here for possible reference later...

// V1 read Request body and conver to a string then unmarshal string into go struct
// body, err := ioutil.ReadAll(req.Body)
// if err != nil {
// 	panic(err)
// }
// log.Println(string(body))
// var newTrainer Test
// err = json.Unmarshal(body, &newTrainer)
// if err != nil {
// 	panic(err)
// }
// log.Println(newTrainer.Name)

// V2 working with parseForm method. NOTE this still works if form is not submitted in typical fashion but through a strigified JSON in ajax request.
// err := req.ParseForm()
// if err != nil {
// 	panic(err)
// }
// fmt.Printf(req.FormValue("name"))
// fmt.Printf(req.FormValue("age"))
// fmt.Printf(req.FormValue("city"))
// fmt.Println(req.Form)

// V3 HOW IT WAS SUPPOSED TO WORK IN THE FIRST PLACE
// var newTrainer Test
// if err := json.NewDecoder(req.Body).Decode(&newTrainer); err != nil {
// 	fmt.Println(err)
// }
// fmt.Println(newTrainer)
// fmt.Printf("%T", newTrainer)
