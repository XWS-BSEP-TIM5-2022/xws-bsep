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

func (store *AuthPostgresStore) Create(authentication *domain.Authentication) (string, error) {
	result := store.db.Create(authentication)
	if result.Error != nil {
		return "error", result.Error
	}
	return "success", nil
}

func (store *AuthPostgresStore) GetAll() (*[]domain.Authentication, error) {
	var auths []domain.Authentication
	result := store.db.Find(&auths)
	if result.Error != nil {
		return nil, result.Error
	}
	return &auths, nil
}

func (store *AuthPostgresStore) Get(id string) (*domain.Authentication, error) {
	authentication := domain.Authentication{}

	result := store.db.Find(&authentication, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &authentication, nil
}

func (store *AuthPostgresStore) DeleteAll() {
	store.db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&domain.Authentication{})
}

func (store *AuthPostgresStore) Insert(auth *domain.Authentication) (string, error) {
	result := store.db.Create(auth)
	if result.Error != nil {
		return "error", result.Error
	}
	return "success", nil
}
