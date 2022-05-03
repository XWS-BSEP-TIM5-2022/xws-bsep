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

// func (service *AuthService) GetAll() (*[]domain.Authentication, error) {
// 	return service.store.GetAll()
// }

func (service *AuthService) Create(auth *domain.Authentication) (string, error) {
	success, err := service.store.Create(auth)
	return success, err
}
