package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func ConnectDB() *mongo.Client {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var data struct {
		MONGO_URL string
	}

	ctx = context.WithValue(ctx, data.MONGO_URL, os.Getenv("MONGO_URL"))
	uri := fmt.Sprintf(`mongodb://%s`,
		ctx.Value(data.MONGO_URL).(string),
	)
	//fmt.Println(uri)

	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(uri))

	err := client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB Employees Database")

	// Creating text indexes for Search Endpoint
	collection := client.Database("employees_db").Collection("employees")
	collection2 := client.Database("employees_db").Collection("personalInfo")

	mod := mongo.IndexModel{
		Keys: bson.D{
			{"firstname", "text"},
			{"lastname", "text"},
			{"desg", "text"},
			{"dpt", "text"},
			{"manager", "text"},
			{"location", "text"},
		},
		// create UniqueIndex option
		//Options: options.Index().SetWeights(bson.D{
		//	{"lastname", 9},
		//	{"firstname", 3},
		//}),
	}

	collection.Indexes().CreateOne(ctx, mod)

	mod = mongo.IndexModel{
		Keys: bson.M{
			"eid": 1, // index in ascendingg order
		},
		// create UniqueIndex option
		Options: options.Index().SetUnique(true),
	}

	collection.Indexes().CreateOne(ctx, mod)
	collection2.Indexes().CreateOne(ctx, mod)

	return client
}
