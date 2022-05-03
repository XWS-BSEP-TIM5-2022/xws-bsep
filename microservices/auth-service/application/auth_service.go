package application

import "auth-service/domain"

type AuthService struct {
	store domain.AuthStore
}

func NewAuthService(store domain.AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (service *AuthService) GetAll() (*[]domain.Authentication, error) {
	return service.store.GetAll()
}

func (service *AuthService) Insert(auth *domain.Authentication) (string, error) {
	success, err := service.store.Insert(auth)
	return success, err
}
