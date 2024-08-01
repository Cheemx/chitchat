package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "chitchatters"
const collName = "userbase"

// MOST IMPORTANT -> this provides you instance of mongo collection -----> 1st step

var UserCollection *mongo.Collection

func ConnectToDB() {
	// client options ---->2nd
	constr := os.Getenv("MONGODB_URI")
	if constr == "" {
		log.Fatalf("Error getting the connection string")
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOption := options.Client().ApplyURI(constr).SetServerAPIOptions(serverAPI)

	// connect to mongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client, err := mongo.Connect(ctx, clientOption)
	defer cancel()

	if err != nil {
		log.Fatalf("error connecting db: %v", err)
	}

	userDatabase := client.Database(dbName)
	UserCollection = userDatabase.Collection(collName)
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	fmt.Println("Collection instance is ready")

}
