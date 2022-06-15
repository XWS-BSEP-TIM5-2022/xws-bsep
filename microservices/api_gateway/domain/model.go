package domain

import (
	"time"
)

type Users struct {
	UsersDetails []User
}

type GenderEnum int8
type RoleEnum int8
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

const (
	Admin RoleEnum = iota
	Registered_User
)

type User struct {
	Id           string
	Name         string
	LastName     string
	MobileNumber string
	Gender       GenderEnum
	Birthday     time.Time
	Email        string
	Biography    string
	Username     string
	Password     string
	IsPublic     bool
	Education    []Education
	Experience   []Experience
	Skills       []Skill
	Interests    []Interest
}

type Education struct {
	Id        string
	Name      string
	Level     EducationEnum
	Place     string
	StartDate time.Time
	EndDate   time.Time
}

type Experience struct {
	Id        string
	Name      string
	Headline  string
	Place     string
	StartDate time.Time
	EndDate   time.Time
}

type Skill struct {
	Id   string
	Name string
}

type Interest struct {
	Id          string
	Name        string
	Description string
}

type Connection struct {
	Id      string
	UserAid string
	UserBid string
}

type Authentication struct {
	Id       string
	Name     string
	Password string
	Role     string
	Date     time.Time
}

type Like struct {
	Id     string
	UserId string
}

type Dislike struct {
	Id     string
	UserId string
}

type Comment struct {
	Id     string
	UserId string
	Text   string
}

type Post struct {
	Id          string
	Text        string
	Image       string
	Links       []string
	DateCreated time.Time
	Likes       []Like
	Dislikes    []Dislike
	Comments    []Comment
	UserId      string
	JobOffer    JobOffer
	Company     Company
	IsJobOffer  bool
}

type PostAgents struct {
	Id          string
	Text        string
	Image       string
	Links       []string
	DateCreated time.Time
	Likes       []Like
	Dislikes    []Dislike
	Comments    []Comment
	UserId      string
	JobOffer    JobOfferAgents
	Company     Company
	IsJobOffer  bool
	ApiToken    string
}

type JobOffer struct {
	Id              string
	Position        Position
	JobDescription  string
	DailyActivities string
	Preconditions   string
}

type JobOfferAgents struct {
	Id              string
	Position        PositionAgents
	JobDescription  string
	DailyActivities string
	Preconditions   string
}

type Position struct {
	Id   string
	Name string
	Pay  float64
}

type PositionAgents struct {
	Id   string
	Name string
	Pay  string
}

type Company struct {
	Id          string
	Name        string
	Description string
	PhoneNumber string
	IsActive    bool
}

type Posts struct {
	AllPosts []Post
}
