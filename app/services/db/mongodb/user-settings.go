package mongodb

import (
	"context"
	"fmt"
	"log"
	"time"

	"polygnosics/app/models"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddSettings() (*mongo.InsertOneResult, error) {
	client, err := Connect()
	if err != nil {
		return nil, err
	}

	userSettings := models.UserSetting{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("user_database").Collection("user_settings")
	insertResult, err := collection.InsertOne(ctx, userSettings)
	if err != nil {
		return nil, err
	}

	if err = client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect mongo client: %s\n", errors.WithStack(err))
	}
	return insertResult, nil
}

func GetSettings(objID *primitive.ObjectID) (*models.UserSetting, error) {
	if objID == nil {
		return nil, fmt.Errorf("Invalid settings ID")
	}

	client, err := Connect()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("user_database").Collection("user_settings")
	result := collection.FindOne(ctx, bson.M{"_id": objID})
	settings := models.UserSetting{}
	if err := result.Decode(&settings); err != nil {
		return nil, err
	}

	if err = client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect mongo client: %s\n", errors.WithStack(err))
	}

	return &settings, nil
}

func DeleteSettings(objID *primitive.ObjectID) (*mongo.DeleteResult, error) {
	if objID == nil {
		return nil, fmt.Errorf("Invalid settings ID")
	}

	client, err := Connect()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	collection := client.Database("user_database").Collection("user_settings")
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return nil, err
	}

	if err = client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect mongo client: %s\n", errors.WithStack(err))
	}
	return result, nil
}
