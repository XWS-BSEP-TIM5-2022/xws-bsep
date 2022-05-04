package domain

import "time"

type Authentication struct {
	Id       string    `gorm:"index:idx_name,unique"`
	Name     string    `gorm:"index:idx_name,unique"`
	Password string    `gorm:"index:idx_name"`
	Role     string    `gorm:"index:idx_name"`
	Date     time.Time `gorm:"index:idx_name"`
}

type RoleEnum int8

const (
	User RoleEnum = iota
	Registered_User
)

func (role RoleEnum) String() string {
	if role == User {
		return "User"
	}
	return "Registered_User"
}
