package persistence

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
	"gorm.io/gorm"
)

type UserPostgresStore struct {
	db *gorm.DB
}

func NewUserPostgresStore(db *gorm.DB) (domain.UserStore, error) {
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		return nil, err
	}
	return &UserPostgresStore{
		db: db,
	}, nil
}

func (store *UserPostgresStore) Get(id string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (store *UserPostgresStore) Update(user *domain.User) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (store *UserPostgresStore) Insert(user *domain.User) (string, error) {
	result := store.db.Create(user)
	if result.Error != nil {
		return "error", result.Error
	}
	return "success", nil
}

func (store *UserPostgresStore) GetAll() (*[]domain.User, error) {
	var users []domain.User
	result := store.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return &users, nil
}

func (store *UserPostgresStore) DeleteAll() {
	store.db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&domain.User{})
}
