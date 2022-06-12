package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	Id               string  `gorm:"index:idx_name,unique"` // id je id usera
	Username         string  `gorm:"index:idx_name,unique"`
	Password         string  `gorm:"index:idx_name"`
	Roles            *[]Role `gorm:"many2many:auth_roles"` //  []*Role - Role             string `gorm:"index:idx_name"`
	VerificationCode string  `gorm:"index:idx_name"`
	ExpirationTime   int64   `gorm:"index:idx_name"`
	APIToken         string  `gorm:"index:idx_name,unique"`
}

type Role struct {
	ID          uint          `gorm:"primaryKey;auto_increment:true"`
	Name        string        `gorm:"index:idx_name,unique"`
	Permissions []*Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
	ID   uint   `gorm:"primaryKey;auto_increment:true"`
	Name string `gorm:"index:idx_name,unique"`
}

func NewAuthCredentials(id, username, password string, roles *[]Role) (*Authentication, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Can not hash password: %w", err)
	}
	credentials := &Authentication{
		Id:               id,
		Username:         username,
		Password:         string(hashedPassword),
		Roles:            roles,
		VerificationCode: "",
		ExpirationTime:   0,
	}
	return credentials, nil
}

func (credentials *Authentication) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password))
	return err == nil
}

func (credentials *Authentication) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Can not hash password: %w", err)
	}
	return string(hashedPassword), nil
}
