package api

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetSingleStatus performs a findOne query provided a task's ObjectId
// represented as a hex string
func GetSingleStatus(database *mongo.Database, id string) (*Status, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.Collection("tasks")
	var status Status
	err = collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// GetAllStatus performs the find all query equivalent to `list(db.tasks.find())`
// in Python/PyMongo
func GetAllStatus(database *mongo.Database) (*[]Status, error) {
	collection := database.Collection("tasks")
	var statusList []Status

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &statusList); err != nil {
		return nil, err
	}

	return &statusList, nil
}

// StopTask marks the requested record as ready to be stopped,
// so the running process will know to shutdown.
func StopTask(database *mongo.Database, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.Collection("tasks")
	updateResult, err := collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"stop_flag": true})
	if err != nil || updateResult.ModifiedCount != 1 {
		return err
	}

	return nil
}
