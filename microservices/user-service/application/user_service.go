package application

import "github.com/sanjadrinic/test_repo/microservices/user_service/domain"

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
