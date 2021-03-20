package db

import (
	"context"
	"log"
	"os"

	"github.com/alexanderschau/ipfs-pinning-service/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collection *mongo.Collection
var Ctx = context.TODO()

func init() {
	env.Load()
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(Ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(Ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	Collection = client.Database(os.Getenv("MONGO_DB")).Collection("pins")
}
