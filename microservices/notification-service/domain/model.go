package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Notification struct {
	Id   primitive.ObjectID `bson:"_id"`
	Time time.Time          `bson:"time"`
}
