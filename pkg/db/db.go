package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeout = 10
const name = "jocastanu"
const username = "example"
const password = "example"
const uri = "mongodb://mongo:27017"

func CreateClient() *mongo.Client {
	log.Printf("Creating client instance for URI: %s", uri)
	cred := options.Credential{Username: username, Password: password}
	clientOptions := options.Client().ApplyURI(uri).SetAuth(cred)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func DisconnectClient(client *mongo.Client) {
	log.Printf("Disconnecting client instance for URI %s", uri)
	ctx, cancel := createTimeoutContext(timeout)
	defer cancel()
	client.Disconnect(*ctx)
}

func createTimeoutContext(tmo int) (*context.Context, context.CancelFunc) {
	log.Printf("Creating timeout context(%d)", tmo)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(tmo)*time.Second)
	return &ctx, cancel
}

func createDeadlineContext(deadline int) (*context.Context, context.CancelFunc) {
	log.Printf("Creating deadline context(%d)", deadline)
	duration := time.Now().Add(time.Duration(deadline) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), duration)
	return &ctx, cancel
}

func PutAny(client *mongo.Client, collectionName string, customStruct interface{}) {
	log.Printf("Putting %+v into collection %s", customStruct, collectionName)
	ctx, cancel := createTimeoutContext(timeout)
	defer cancel()

	database := client.Database(name)
	recordCollection := database.Collection(collectionName)
	_, err := recordCollection.InsertOne(*ctx, customStruct)
	if err != nil {
		log.Fatal(err)
	}
}

func GetAny(client *mongo.Client, collectionName string, key string, value interface{}) []primitive.M {
	database := client.Database(name)
	collection := database.Collection(collectionName)
	ctx, cancel := createDeadlineContext(timeout)
	defer cancel()

	bsonFilter := bson.M{}
	if key != "" && value != nil {
		bsonFilter = bson.M{key: value}
		log.Printf("Get all elements from collection %s for query %s:%s", collectionName, key, value)
	} else {
		log.Printf("Get all elements from collection %s", collectionName)
	}

	cur, err := collection.Find(*ctx, bsonFilter)
	if err != nil {
		log.Fatal(err)
	}
	var filtered []bson.M
	if err = cur.All(*ctx, &filtered); err != nil {
		log.Fatal(err)
	}
	return filtered
}

func GetOne(client *mongo.Client, collectionName string, key string, value interface{}) *mongo.SingleResult {
	log.Printf("Get one element from collection %s for query %s:%s", collectionName, key, value)
	database := client.Database(name)
	collection := database.Collection(collectionName)
	ctx, cancel := createDeadlineContext(timeout)
	defer cancel()

	result := collection.FindOne(*ctx, bson.M{key: value})

	if result.Err() != nil {
		log.Printf("No elements found in collection %s for query %s:%s: %s", collectionName, key, value, result.Err().Error())
	}
	return result
}

func DeleteAll(client *mongo.Client, collectionName string) {
	database := client.Database(name)
	collection := database.Collection(collectionName)
	ctx, cancel := createDeadlineContext(timeout)
	defer cancel()

	bsonFilter := bson.M{}
	result, err := collection.DeleteMany(*ctx, bsonFilter)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Deleted %d elements from %s collection", result.DeletedCount, collectionName)
}

func CountDocuments(client *mongo.Client, collectionName string) int {
	database := client.Database(name)
	collection := database.Collection(collectionName)
	ctx, cancel := createDeadlineContext(timeout)
	defer cancel()

	count, err := collection.CountDocuments(*ctx, bson.M{}, nil)
	if err != nil {
		log.Fatal(err)
	}
	return int(count)
}
