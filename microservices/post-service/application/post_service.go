package application

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostService struct {
	store      domain.PostStore
	eventStore domain.EventStore
}

func NewPostService(store domain.PostStore, eventStore domain.EventStore) *PostService {
	return &PostService{
		store:      store,
		eventStore: eventStore,
	}
}

func (service *PostService) Get(id primitive.ObjectID) (*domain.Post, error) {
	return service.store.Get(id)
}

func (service *PostService) GetAll() ([]*domain.Post, error) {
	return service.store.GetAll()
}

func (service *PostService) Insert(post *domain.Post) (string, error) {
	success, err := service.store.Insert(post)
	return success, err
}

func (service *PostService) Update(post *domain.Post) (string, error) {
	success, err := service.store.Update(post)
	return success, err
}

func (service *PostService) GetAllByUser(id string) ([]*domain.Post, error) {
	return service.store.GetAllByUser(id)
}

func (service *PostService) LikePost(post *domain.Post, id string) (string, error) {
	return service.store.LikePost(post, id)
}

func (service *PostService) DislikePost(post *domain.Post, id string) (string, error) {
	return service.store.DislikePost(post, id)
}

func (service *PostService) CommentPost(post *domain.Post, id string, text string) (string, error) {
	return service.store.CommentPost(post, id, text)
}

func (service *PostService) UpdateCompanyInfo(company *domain.Company, oldName string) (string, error) {
	return service.store.UpdateCompanyInfo(company, oldName)
}

func (service *PostService) NewEvent(event *domain.Event) (*domain.Event, error) {
	_, err := service.eventStore.NewEvent(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (service *PostService) GetAllEvents() ([]*domain.Event, error) {
	return service.eventStore.GetAllEvents()
}
