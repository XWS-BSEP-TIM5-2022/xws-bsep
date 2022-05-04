package domain

import "time"

type Authentication struct {
	Id       string    `gorm:"index:idx_name,unique"`
	Name     string    `gorm:"index:idx_name,unique"`
	Password string    `gorm:"index:idx_name"`
	Role     string    `gorm:"index:idx_name"` // TODO Sanja: ispraviti - enum
	Date     time.Time `gorm:"index:idx_name"`
}
