package domain

import "time"

type Users struct {
	UsersDetails []User
}

type GenderEnum int8

const (
	Female GenderEnum = iota
	Male
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
