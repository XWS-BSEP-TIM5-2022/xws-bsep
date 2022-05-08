package persistence

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	"gorm.io/gorm"
)

type AuthPostgresStore struct {
	db *gorm.DB
}

func NewAuthPostgresStore(db *gorm.DB) (*AuthPostgresStore, error) {
	err := db.AutoMigrate(&domain.Authentication{})
	if err != nil {
		return nil, err
	}
	return &AuthPostgresStore{
		db: db,
	}, nil
}

func (store *AuthPostgresStore) Create(auth *domain.Authentication) (*domain.Authentication, error) {
	err := store.db.Create(auth)
	return auth, err.Error
}

func (store *AuthPostgresStore) FindByUsername(username string) (*domain.Authentication, error) {
	var auth domain.Authentication
	err := store.db.First(&auth, "username = ?", username)
	return &auth, err.Error
}

func (store *AuthPostgresStore) FindAll() (*[]domain.Authentication, error) {
	var auths []domain.Authentication
	result := store.db.Find(&auths)
	if result.Error != nil {
		return nil, result.Error
	}
	return &auths, nil
}

func (store *AuthPostgresStore) DeleteAll() {
	store.db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&domain.Authentication{})
}

func (store *AuthPostgresStore) Insert(auth *domain.Authentication) error {
	result := store.db.Create(auth)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
