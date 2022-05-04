package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	LastName     string             `bson:"last_name"`
	MobileNumber string             `bson:"mobile_number"`
	Gender       GenderEnum         `bson:"gender"`
	Birthday     time.Time          `bson:"birthday"`
	Email        string             `bson:"email"`
	Biography    string             `bson:"biography"`
	Username     string             `bson:"username"`
	Password     string             `bson:"password"`
}

type GenderEnum int8

const (
	Female GenderEnum = iota
	Male
)
