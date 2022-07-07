package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PostsID struct {
	Id string
}

type Posts struct {
	AllPosts []Post
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

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name" validate:"required"`
	LastName     string             `bson:"last_name" validate:"required"`
	MobileNumber string             `bson:"mobile_number"`
	Username     string             `bson:"username"`
	Gender       GenderEnum         `bson:"gender"`
	Birthday     time.Time          `bson:"birthday" validate:"required"`
	Email        string             `bson:"email" validate:"required"`
	Biography    string             `bson:"biography"`
	IsPublic     bool               `bson:"is_public"`
	Education    []Education        `bson:"education"`
	Experience   []Experience       `bson:"experience"`
	Skills       []Skill            `bson:"skills"`
	Interests    []Interest         `bson:"interests"`
	IsActive     bool               `bson:"is_active"`
	Role         []string           `bson:"role"`
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
