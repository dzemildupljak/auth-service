package mngdb

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DbConnection(ctx context.Context) *mongo.Client {
	mongouri := os.Getenv("MONGODB_LOCAL_URI")

	mongoconn := options.Client().ApplyURI(mongouri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")
	return mongoclient

	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mngroot:mongodb@127.0.0.1:27017/auth_service_db"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // TODO add ctx WithTimeout in main.go
	// // ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// err = client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer client.Disconnect(ctx)

	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("\n\n\n", databases)
}
