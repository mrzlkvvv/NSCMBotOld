package database

import (
	"context"
	"log"

	"github.com/KirillMerz/NSCMTelegramBot/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (db *Database) IsUserRegistered(UserID int64) (bool, error) {
	var user models.User

	filter := bson.D{{Key: "_id", Value: UserID}}
	err := db.users.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (db *Database) RegisterUser(user models.User) error {
	_, err := db.users.InsertOne(context.TODO(), user)
	return err
}

func (db *Database) GetUserByID(UserID int64) (models.User, error) {
	var user models.User

	filter := bson.D{{Key: "_id", Value: UserID}}
	err := db.users.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (db *Database) UnregisterUser(UserID int64) error {
	filter := bson.D{{Key: "_id", Value: UserID}}
	_, err := db.users.DeleteOne(context.TODO(), filter)
	return err
}

func (db *Database) UpdateResults(UserID int64, results models.Results) error {
	_, err := db.results.UpdateByID(context.TODO(), UserID, results)
	return err
}
