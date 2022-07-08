package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	Id       primitive.ObjectID `bson:"_id"`
	Receiver string             `bson:"receiver"`
	Content  string             `bson:"content"`
	Time     time.Time          `bson:"time"`
}

type Conversation struct {
	Id       primitive.ObjectID `bson:"_id"`
	User1    string             `bson:"user1"`
	User2    string             `bson:"user2"`
	Messages []Message          `bson:"messages"`
}
