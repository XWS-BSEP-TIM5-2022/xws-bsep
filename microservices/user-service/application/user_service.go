package application

import (
	"context"
	"errors"
	"net/mail"
	"unicode"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	store        domain.UserStore
	orchestrator *CreateUserOrchestrator
	eventStore   domain.EventStore
}

func NewUserService(store domain.UserStore, orchestrator *CreateUserOrchestrator, eventStore domain.EventStore) *UserService {
	return &UserService{
		store:        store,
		orchestrator: orchestrator,
		eventStore:   eventStore,
	}
}

func (service *UserService) GetAll(ctx context.Context) ([]*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetAll(ctx)
}

func (service *UserService) GetAllPublic(ctx context.Context) ([]*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetAllPublic(ctx)
}

func (service *UserService) Insert(user *domain.User) (*domain.User, error) {
	_, err := service.store.Insert(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) Update(ctx context.Context, user *domain.User) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "Update service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	success, err := service.store.Update(ctx, user)
	return success, err
}

func (service *UserService) UpdatePostNotification(ctx context.Context, user *domain.User) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "Update service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	success, err := service.store.Update(ctx, user)
	return success, err
}

func (service *UserService) UpdateBasicInfo(ctx context.Context, user *domain.User) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateBasicInfo service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	success, err := service.store.UpdateBasicInfo(ctx, user)
	return success, err
}

func (service *UserService) UpdatePrivacy(ctx context.Context, user *domain.User) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdatePrivacyInfo service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	success, err := service.store.UpdatePrivacy(ctx, user)
	return success, err
}

func (service *UserService) UpdateExperienceAndEducation(ctx context.Context, user *domain.User) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateExperienceAndEducation service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	success, err := service.store.UpdateExperienceAndEducation(ctx, user)
	return success, err
}

func (service *UserService) UpdateSkillsAndInterests(ctx context.Context, user *domain.User) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateSkillsAndInterests service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	success, err := service.store.UpdateSkillsAndInterests(ctx, user)
	return success, err
}

func (service *UserService) Get(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "Get service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.Get(ctx, id)
}

func (service *UserService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetByUsername service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetByUsername(ctx, username)
}

func (service *UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetByEmail service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetByEmail(email)
}

func (service *UserService) GetById(ctx context.Context, userId string) (*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "GetById service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetById(ctx, userId)
}

func (service *UserService) Search(ctx context.Context, criteria string) ([]*domain.User, error) {
	span := tracer.StartSpanFromContext(ctx, "Search service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.Search(ctx, criteria)
}

func (service *UserService) UpdateIsActiveById(ctx context.Context, userId string) error {
	span := tracer.StartSpanFromContext(ctx, "UpdateIsActiveById service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.UpdateIsActiveById(ctx, userId)
}

func (service *UserService) GetIdByEmail(ctx context.Context, email string) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "GetIdByEmail service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)
	return service.store.GetIdByEmail(ctx, email)
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

func (service *UserService) NewEvent(event *domain.Event) (*domain.Event, error) {
	_, err := service.eventStore.NewEvent(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (service *UserService) GetAllEvents() ([]*domain.Event, error) {
	return service.eventStore.GetAllEvents()
}
