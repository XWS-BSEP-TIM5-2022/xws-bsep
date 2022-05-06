package application

import "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"

type AuthService struct {
	store domain.AuthStore
}

func NewAuthService(store domain.AuthStore) *AuthService {
	return &AuthService{
		store: store,
	}
}

func (service *AuthService) FindByUsername(username string) (*domain.Authentication, error) {
	success, err := service.store.FindByUsername(username)
	return success, err
}

func (service *AuthService) Create(auth *domain.Authentication) (*domain.Authentication, error) {
	success, err := service.store.Create(auth)
	return success, err
}
