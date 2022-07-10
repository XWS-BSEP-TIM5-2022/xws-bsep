package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Notification struct {
	Id     primitive.ObjectID   `bson:"_id"`
	Date   time.Time            `bson:"date"`
	Text   string               `bson:"text"`
	UserId string               `bson:"user_id"`
	Type   NotificationTypeEnum `bson:"type"`
	Read   bool                 `bson:"read"`
}

type NotificationTypeEnum int8

const (
	Message NotificationTypeEnum = iota
	Follow
	Post
)

type Event struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId string             `bson:"user_id"`
	Text   string             `bson:"action"`
	Date   time.Time          `bson:"date"`
}
