package main

import (
	"context"
	"encoding/json"
	"fmt"
	// "log"
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

func QueryDB(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5174")

	state_prefix := req.URL.Query().Get("state_prefix")

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
