package domain

import "time"

type User struct {
	// Id   int `bson:"_id"`
	// Name string             `bson:"name"`
	// Id   int    `json:"id"`
	// Name string `json:"name"`
	Id           string     `gorm:"index:idx_name,unique"`
	Name         string     `gorm:"index:idx_name"`
	LastName     string     `gorm:"index:idx_name"`
	MobileNumber string     `gorm:"index:idx_name"`
	Gender       GenderEnum `gorm:"index:idx_name"`
	Birthday     time.Time  `gorm:"index:idx_name"`
	Email        string     `gorm:"index:idx_name,unique"`
	Biography    string     `gorm:"index:idx_name"`
	Username     string     `gorm:"index:idx_name,unique"`
	Password     string     `gorm:"index:idx_name"`
}

type GenderEnum int8

const (
	Female GenderEnum = iota
	Male
)

func (gender GenderEnum) String() string {
	if gender == Female {
		return "Female"
	}
	return "Male"
}
