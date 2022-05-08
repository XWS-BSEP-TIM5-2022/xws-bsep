package domain

import (
	"time"
)

type Users struct {
	UsersDetails []User
}

type GenderEnum int8
type RoleEnum int8

const (
	Female GenderEnum = iota
	Male
)

const (
	Admin RoleEnum = iota
	Registered_User
)

type User struct {
	Id           string
	Name         string
	LastName     string
	Email        string
	MobileNumber string
	Gender       GenderEnum
	Birthday     time.Time
	Username     string
	Biography    string
	Password     string
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
	Images      []string
	Links       []string
	DateCreated time.Time
	Likes       []Like
	Dislikes    []Dislike
	Comments    []Comment
	UserId      string
}

type Posts struct {
	AllPosts []Post
}
