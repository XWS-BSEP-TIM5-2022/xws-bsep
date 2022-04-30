package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	// Id   int `bson:"_id"`
	// Name string             `bson:"name"`
	// Id   int    `json:"id"`
	// Name string `json:"name"`
	Id   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}
