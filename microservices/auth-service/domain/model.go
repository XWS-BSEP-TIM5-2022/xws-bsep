package domain

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	Id       string `gorm:"index:idx_name,unique"` // id je id usera
	Username string `gorm:"index:idx_name,unique"`
	Password string `gorm:"index:idx_name"`
	Role     string `gorm:"index:idx_name"`
}

func NewAuthCredentials(id, username, password, role string) (*Authentication, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}
	credentials := &Authentication{
		Id:       id,
		Username: username,
		Password: string(hashedPassword),
		Role:     role, // TODO: modifikovati kada dodamo role
	}
	return credentials, nil
}

func (credentials *Authentication) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(credentials.Password), []byte(password))
	return err == nil
}
