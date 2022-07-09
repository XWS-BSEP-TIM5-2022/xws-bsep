package application

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/notification_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationService struct {
	store domain.NotificationStore
}

func NewNotificationService(store domain.NotificationStore) *NotificationService {
	return &NotificationService{
		store: store,
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
