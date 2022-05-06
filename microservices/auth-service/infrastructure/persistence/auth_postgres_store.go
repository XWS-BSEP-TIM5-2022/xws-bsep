package persistence

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	"gorm.io/gorm"
)

type AuthPostgresStore struct {
	db *gorm.DB
}

func NewAuthPostgresStore(db *gorm.DB) (domain.AuthStore, error) {
	err := db.AutoMigrate(&domain.Authentication{})
	if err != nil {
		return nil, err
	}
	return &AuthPostgresStore{
		db: db,
	}, nil
}

func (store *AuthPostgresStore) Create(auth *domain.Authentication) (*domain.Authentication, error) {
	result := store.db.Create(auth)

	return auth, result.Error
}

func (store *AuthPostgresStore) FindByUsername(username string) (*domain.Authentication, error) {
	var authCredentials domain.Authentication
	result := store.db.First(&authCredentials, "username = ?", username)

	return &authCredentials, result.Error
}
