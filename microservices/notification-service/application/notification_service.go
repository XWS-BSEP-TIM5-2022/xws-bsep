package application

import (
	"context"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
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

func (service *NotificationService) GetById(ctx context.Context, id primitive.ObjectID) (*domain.Notification, error) {
	span := tracer.StartSpanFromContext(ctx, "GetById service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetById(ctx, id)
}

func (service *NotificationService) GetAll(ctx context.Context) ([]*domain.Notification, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAll(ctx)
}

func (service *NotificationService) Insert(ctx context.Context, notification *domain.Notification) (string, error) {
	span := tracer.StartSpanFromContext(ctx, "Insert service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	success, err := service.store.Insert(ctx, notification)
	return success, err
}

func (service *NotificationService) GetAllByUser(ctx context.Context, id string) ([]*domain.Notification, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllByUser service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllByUser(ctx, id)
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
