package application

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/notification_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationService struct {
	store      domain.NotificationStore
	eventStore domain.EventStore
}

func NewNotificationService(store domain.NotificationStore, eventStore domain.EventStore) *NotificationService {
	return &NotificationService{
		store:      store,
		eventStore: eventStore,
	}
}

func (service *NotificationService) GetById(id primitive.ObjectID) (*domain.Notification, error) {
	return service.store.GetById(id)
}

func (service *NotificationService) GetAll() ([]*domain.Notification, error) {
	return service.store.GetAll()
}

func (service *NotificationService) Insert(notification *domain.Notification) (string, error) {
	success, err := service.store.Insert(notification)
	return success, err
}

func (service *NotificationService) GetAllByUser(id string) ([]*domain.Notification, error) {
	return service.store.GetAllByUser(id)
}

func (service *NotificationService) NewEvent(event *domain.Event) (*domain.Event, error) {
	_, err := service.eventStore.NewEvent(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (service *NotificationService) GetAllEvents() ([]*domain.Event, error) {
	return service.eventStore.GetAllEvents()
}
