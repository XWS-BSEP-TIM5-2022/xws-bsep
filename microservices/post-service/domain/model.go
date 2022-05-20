package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Like struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId string             `bson:"user_id" validate:"required"`
}

type Dislike struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId string             `bson:"user_id" validate:"required"`
}

type Comment struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId string             `bson:"user_id" validate:"required"`
	Text   string             `bson:"text"`
}

type Post struct {
	Id          primitive.ObjectID `bson:"_id"`
	Text        string             `bson:"text"`
	Images      []string           `bson:"images"`
	Links       []string           `bson:"links"`
	DateCreated time.Time          `bson:"date_created" validate:"required"` // datetime
	Likes       []Like             `bson:"likes"`
	Dislikes    []Dislike          `bson:"dislikes"`
	Comments    []Comment          `bson:"comments"`
	UserId      string             `bson:"user_id"`
}
