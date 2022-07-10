package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserConn struct {
	UserID   string
	IsPublic bool
}

type Event struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId string             `bson:"user_id"`
	Text   string             `bson:"action"`
	Date   time.Time          `bson:"date"`
}
