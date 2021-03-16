package db

import (
	"context"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	return &ctx, cancel
}
