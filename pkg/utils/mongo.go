package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertToObjectId(id string) primitive.ObjectID {
	objID, _ := primitive.ObjectIDFromHex(id)
	return objID
}
