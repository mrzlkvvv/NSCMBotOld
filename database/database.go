package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/KirillMerz/NSCMTelegramBot/models"
)

const DATABASE_NAME = "NSCMTelegramBot"

type Database struct {
	users   *mongo.Collection
	results *mongo.Collection
}

func New(uri string) *Database {
	log.Println("DB: connection started...")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatalln("DB: connection failed:", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalln("DB: ping failed:", err)
	}

	db := client.Database(DATABASE_NAME)

	log.Println("DB: succesfully connected!")

	return &Database{
		users:   db.Collection("users"),
		results: db.Collection("results"),
	}
}

func (db *Database) RegisterUser(user models.User, results models.Results) error {
	_, err := db.users.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}

	_, err = db.results.InsertOne(context.TODO(), results)

	return err
}

func (db *Database) UnregisterUser(UserID int64) error {
	filter := bson.D{{Key: "_id", Value: UserID}}

	_, err := db.users.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	_, err = db.results.DeleteOne(context.TODO(), filter)

	return err
}

func (db *Database) GetResults(UserID int64) (models.Results, error) {
	filter := bson.D{{Key: "_id", Value: UserID}}

	var results models.Results
	err := db.results.FindOne(context.TODO(), filter).Decode(&results)

	return results, err
}

func (db *Database) ReplaceResults(UserID int64, results models.Results) error {
	filter := bson.D{{Key: "_id", Value: UserID}}
	_, err := db.results.ReplaceOne(context.TODO(), filter, results)
	return err
}

func (db *Database) GetAllUsers() ([]models.User, error) {
	cur, err := db.users.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = cur.All(context.TODO(), &users)

	return users, err
}
