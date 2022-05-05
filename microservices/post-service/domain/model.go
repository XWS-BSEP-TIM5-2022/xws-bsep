package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Like struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId string             `bson:"user_id"`
	// datum i vreme ?
}

type Dislike struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId string             `bson:"user_id"`
}

type Comment struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId string             `bson:"user_id"`
	Text   string             `bson:"text"`
}

type Post struct {
	Id          primitive.ObjectID `bson:"_id"`
	Text        string             `bson:"text"`
	Images      []string           `bson:"images"`
	Links       []string           `bson:"links"`
	DateCreated time.Time          `bson:"date_created"`
	Likes       []Like             `bson:"likes"`
	Dislikes    []Dislike          `bson:"dislikes"`
	Comments    []Comment          `bson:"comments"`
	UserId      string             `bson:"user_id"`
}
