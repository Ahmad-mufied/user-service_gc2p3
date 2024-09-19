package utils

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	//dbName := config.Viper.GetString("MONGO_DB")
	return client.Database("students").Collection(collectionName)
}

// CheckDocumentExists checks if a document exists in the specified collection by its _id
func CheckDocumentExists(ctx context.Context, collection *mongo.Collection, id primitive.ObjectID) (bool, error) {
	count, err := collection.CountDocuments(ctx, bson.M{"_id": id}, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
