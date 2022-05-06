package application

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	store domain.UserStore
}

func NewUserService(store domain.UserStore) *UserService {
	return &UserService{
		store: store,
	}
}

func (service *UserService) GetAll() ([]*domain.User, error) {
	return service.store.GetAll()
}

func (service *UserService) GetAllPublic() ([]*domain.User, error) {
	return service.store.GetAllPublic()
}

func (service *UserService) Insert(user *domain.User) (*domain.User, error) {
	_, err := service.store.Insert(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) Update(user *domain.User) (string, error) {
	success, err := service.store.Update(user)
	return success, err
}

func (service *UserService) Get(id primitive.ObjectID) (*domain.User, error) {
	return service.store.Get(id)
}

func (service *UserService) GetByUsername(username string) (*domain.User, error) {
	return service.store.GetByUsername(username)
}

func (service *UserService) GetByEmail(email string) (*domain.User, error) {
	return service.store.GetByEmail(email)
}

func (service *UserService) GetById(userId string) (*domain.User, error) {
	return service.store.GetById(userId)
}
func (service *UserService) Search(criteria string) ([]*domain.User, error) {
	return service.store.Search(criteria)
}
