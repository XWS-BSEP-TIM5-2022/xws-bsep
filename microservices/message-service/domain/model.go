package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	Id      primitive.ObjectID `bson:"_id"`
	Sender  primitive.ObjectID `bson:"senderId"`
	Content string             `bson:"content"`
	Time    time.Time          `bson:"time"`
}

type Conversation struct {
	Id       primitive.ObjectID `bson:"_id"`
	User1    primitive.ObjectID `bson:"user1"`
	User2    primitive.ObjectID `bson:"user2"`
	Messages []Message          `bson:"messages"`
}
