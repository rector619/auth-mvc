package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//DBinstance func
func DBinstance() *mongo.Client {
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatal("Error loading .env file")
    }

    MongoDb := os.Getenv("MONGODB_URL")
    ctx, cancel := context.WithTimeout(context.Background(), 50*time.Minute)
    defer cancel()
 

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDb))
    if err != nil {
        log.Fatal(err)
    }

    if err := client.Ping(ctx, readpref.Primary()); err != nil {
        fmt.Println(err)
        fmt.Println("error pinging mongo")
    }
   
    fmt.Println("Connected to MongoDB!")

    return client
}

//Client Database instance
var Client *mongo.Client = DBinstance()

//OpenCollection is a  function makes a connection with a collection in the database
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
    DbName := os.Getenv("DB_NAME")

    var collection *mongo.Collection = client.Database(DbName).Collection(collectionName)

    return collection
}
