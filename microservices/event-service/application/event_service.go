package application

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/event_service/domain"
)

type EventService struct {
	store domain.EventStore
}

func NewEventService(store domain.EventStore) *EventService {
	return &EventService{
		store: store,
	}
}

func (service *EventService) GetAllEvents() ([]*domain.Event, error) {
	return service.store.GetAllEvents()
}
