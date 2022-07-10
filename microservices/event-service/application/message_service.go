package application

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/event_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventService struct {
	store domain.EventStore
}

func NewEventService(store domain.EventStore) *EventService {
	return &EventService{
		store: store,
	}
}

func (service *EventService) GetById(id primitive.ObjectID) (*domain.Event, error) {
	return service.store.GetById(id)
}
