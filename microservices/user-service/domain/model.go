package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	Status       UserStatus         `bson:"status"` // SD: SAGA
	CreatedAt    time.Time          `bson:"created_at"`
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

type UserStatus int8

const (
	PendingApproval UserStatus = iota
	Approved
	Cancelled
)

func (status UserStatus) String() string {
	switch status {
	case PendingApproval:
		return "Pending Approval"
	case Approved:
		return "Approved"
	case Cancelled:
		return "Cancelled"
	}
	return "Unknown"
}
