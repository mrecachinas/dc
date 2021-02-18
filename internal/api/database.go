package api

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetAllStatus(database *mongo.Database) (*[]Status, error) {
	collection := database.Collection("tasks")
	var statusList []Status
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &statusList)
	if err != nil {
		return nil, err
	}

	return &statusList, nil
}

func StopTask(database *mongo.Database, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	collection := database.Collection("tasks")
	updateResult, err := collection.UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"stop_flag": true})
	if err != nil || updateResult.ModifiedCount != 1 {
		return err
	}

	return nil
}
