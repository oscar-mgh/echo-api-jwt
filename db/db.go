package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient = Connect()
var DBName = "echo_api"
var CategoryCollection = MongoClient.Database(DBName).Collection("categories")
var CourseCollection = MongoClient.Database(DBName).Collection("courses")
var CourseThumbnailCollection = MongoClient.Database(DBName).Collection("thumbnails")
var UsersCollection = MongoClient.Database(DBName).Collection("users")
var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017/" + DBName)

func Connect() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
		return client
	}
	log.Println("Connected to mongo")
	return client
}

func TestConnection() int {
	err := MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		return 0
	}
	return 1
}
