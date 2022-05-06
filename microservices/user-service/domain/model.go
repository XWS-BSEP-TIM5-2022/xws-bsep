package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	// Password     string             `bson:"password"`
	IsPublic   bool         `bson:"is_public"`
	Education  []Education  `bson:"education"`
	Experience []Experience `bson:"experience"`
	Skills     []Skill      `bson:"skills"`
	Interests  []Interest   `bson:"interests"`
}

type Education struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Level     EducationEnum      `bson:"level"`
	Place     string             `bson:"place"`
	StartDate time.Time          `bson:"start_date"`
	EndDate   time.Time          `bson:"end_date"`
}

type Experience struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Headline  string             `bson:"headline"`
	Place     string             `bson:"place"`
	StartDate time.Time          `bson:"start_date"`
	EndDate   time.Time          `bson:"end_date"`
}

type Skill struct {
	Id   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type Interest struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
}

type GenderEnum int8
type EducationEnum int8

const (
	Female GenderEnum = iota
	Male
)

const (
	Primary EducationEnum = iota
	Secondary
	Bachelor
	Master
	Doctorate
)
