package application

import (
	"errors"
	"net/mail"
	"unicode"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	store        domain.UserStore
	orchestrator *CreateUserOrchestrator
}

func NewUserService(store domain.UserStore, orchestrator *CreateUserOrchestrator) *UserService {
	return &UserService{
		store:        store,
		orchestrator: orchestrator,
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

func (service *UserService) UpdateBasicInfo(user *domain.User) (string, error) {
	success, err := service.store.UpdateBasicInfo(user)
	return success, err
}

func (service *UserService) UpdateExperienceAndEducation(user *domain.User) (string, error) {
	success, err := service.store.UpdateExperienceAndEducation(user)
	return success, err
}

func (service *UserService) UpdateSkillsAndInterests(user *domain.User) (string, error) {
	success, err := service.store.UpdateSkillsAndInterests(user)
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

func (service *UserService) UpdateIsActiveById(userId string) error {
	return service.store.UpdateIsActiveById(userId)
}

func (service *UserService) GetIdByEmail(email string) (string, error) {
	return service.store.GetIdByEmail(email)
}

func (service *UserService) Create(user *domain.User, username, password string) error {
	userDetails := mapNewUser(user, username, password)
	err := service.orchestrator.Start(userDetails)
	if err != nil {
		return err
	}
	return nil
}

func (service *UserService) Delete(user *domain.User) error {
	return service.store.DeleteUser(user.Id.Hex(), user.Email)
}

func checkUsernameCriteria(username string) error {
	if len(username) == 0 {
		return errors.New("Username should not be empty")
	}
	for _, char := range username {
		if unicode.IsSpace(int32(char)) {
			return errors.New("Username should not contain any spaces")
		}
	}
	return nil
}

func (service *UserService) CheckEmailCriteria(email string) error {
	if len(email) == 0 {
		return errors.New("Email should not be empty")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("Email is invalid.")
	}
	return nil
}
