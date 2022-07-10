package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Like struct {
	Id     primitive.ObjectID `bson:"_id" validate:"required"`
	UserId string             `bson:"user_id" validate:"required"`
}

type Dislike struct {
	Id     primitive.ObjectID `bson:"_id" validate:"required"`
	UserId string             `bson:"user_id" validate:"required"`
}

type Comment struct {
	Id     primitive.ObjectID `bson:"_id" validate:"required"`
	UserId string             `bson:"user_id" validate:"required"`
	Text   string             `bson:"text" validate:"required"`
}

type Post struct {
	Id          primitive.ObjectID `bson:"_id" validate:"required"`
	Text        string             `bson:"text" validate:"required"`
	Image       string             `bson:"image"`
	Links       []string           `bson:"links"`
	DateCreated time.Time          `bson:"date_created" validate:"required"`
	Likes       []Like             `bson:"likes"`
	Dislikes    []Dislike          `bson:"dislikes"`
	Comments    []Comment          `bson:"comments"`
	UserId      string             `bson:"user_id" validate:"required"`
	JobOffer    JobOffer           `bson:"job_offer"`
	Company     Company            `bson:"company"`
	IsJobOffer  bool               `bson:"is_job_offer"`
}

type JobOffer struct {
	Id              primitive.ObjectID `bson:"_id"`
	Position        Position           `bson:"position"`
	JobDescription  string             `bson:"job_description"`
	DailyActivities string             `bson:"daily_activities"`
	Preconditions   string             `bson:"preconditions"`
}

type Company struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	PhoneNumber string             `bson:"phone_number"`
	IsActive    bool               `bson:"is_active"`
}

type Position struct {
	Id   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
	Pay  float64            `bson:"pay"`
}

type Event struct {
	Id     primitive.ObjectID `bson:"_id"`
	UserId string             `bson:"user_id"`
	Text   string             `bson:"action"`
	Date   time.Time          `bson:"date"`
}
