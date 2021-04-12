package db

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const timeout = 10

func CreateClient() *mongo.Client {
	cred := options.Credential{Username: "example", Password: "example"}
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017").SetAuth(cred)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func DisconnectClient(client *mongo.Client) {
	ctx, cancel := createTimeoutContext(timeout)
	defer cancel()
	client.Disconnect(*ctx)
}

func GetDatabaseNames(client *mongo.Client) string {
	ctx, cancel := createTimeoutContext(timeout)
	defer cancel()
	databases, err := client.ListDatabaseNames(*ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	return strings.Join(databases, " ")
}

func createTimeoutContext(tmo int) (*context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(tmo)*time.Second)
	return &ctx, cancel
}

func createDeadlineContext(deadline int) (*context.Context, context.CancelFunc) {
	duration := time.Now().Add(time.Duration(deadline) * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), duration)
	return &ctx, cancel
}

func PutAny(client *mongo.Client, customStruct interface{}) {
	ctx, cancel := createTimeoutContext(timeout)
	defer cancel()

	database := client.Database("jocastanu")
	recordCollection := database.Collection("records")
	result, err := recordCollection.InsertOne(*ctx, customStruct)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.InsertedID)

}

func GetAny(client *mongo.Client, collectionName string, key string, value interface{}) []primitive.M {
	database := client.Database("jocastanu")
	collection := database.Collection(collectionName)
	ctx, cancel := createDeadlineContext(timeout)
	defer cancel()

	cur, err := collection.Find(*ctx, bson.M{key: value})
	if err != nil {
		log.Fatal(err)
	}
	var filtered []bson.M
	if err = cur.All(*ctx, &filtered); err != nil {
		log.Fatal(err)
	}
	return filtered
}
