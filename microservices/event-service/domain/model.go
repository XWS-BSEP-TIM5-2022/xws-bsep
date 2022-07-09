package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	Id   primitive.ObjectID `bson:"_id"`
	Date time.Time          `bson:"date"`
}
