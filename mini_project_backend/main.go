package main

import (
	"context"
	"encoding/json"
	"fmt"
    "net/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type State struct {
	ID				primitive.ObjectID `bson:"_id"`
	State			string
	Abbreviation 	string
}

const uri = "mongodb://localhost:27017"
var client *mongo.Client

// Initializes MongoDB client.
func InitClient()(*mongo.Client) {
	client, err := mongo.Connect(context.TODO(), options.Client().
		ApplyURI(uri))
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Connected to %s\n", uri)
	}
	return client
}

// GET request handler for /query
// Params
//		* state_prefix: a string representing the State name's prefix
// Returns
// 		JSON containing a list of State names that's starting with `state_prefix` (case insensitive).
//		The list will be sorted in ascending order.
//		When `state_prefix` is an empty string, returns all possible States in the U.S.
func QueryDB(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
	// The following header needs to be set. For more details, check the following post:
	// https://stackoverflow.com/questions/61238680/access-to-fetch-at-from-origin-http-localhost3000-has-been-blocked-by-cors
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5174")

	state_prefix := req.URL.Query().Get("state_prefix")

	// Query db and get State names starting with `state_prefix`
	collection := client.Database("project_dataset").Collection("us_states")
	filter := bson.D{{"state", bson.D{{"$regex", fmt.Sprintf("^%s", state_prefix)}, {"$options", "i"}}}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var states []State
	if err = cursor.All(context.TODO(), &states); err != nil {
		panic(err)
	}

	// Get State names from the db returned messages
	var state_names []string
	for _, state := range states {
		cursor.Decode(&state)
		state_names = append(state_names, state.State)
	}

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(state_names)
}

func main() {
	client = InitClient()
	defer client.Disconnect(context.Background())

    http.HandleFunc("/query", QueryDB)

    http.ListenAndServe(":3000", nil)
}
