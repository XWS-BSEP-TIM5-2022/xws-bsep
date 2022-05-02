package application

import "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"

type UserService struct {
	store domain.UserStore
}

func NewUserService(store domain.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) GetAll() (*[]domain.User, error) {
	return service.store.GetAll()
}

func (service *UserService) Insert(user *domain.User) (string, error) {
	success, err := service.store.Insert(user)
	return success, err
}
