package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	Id primitive.ObjectID `bson:"_id"`
}
