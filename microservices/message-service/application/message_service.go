package application

import (
	"context"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService struct {
	store      domain.MessageStore
	eventStore domain.EventStore
}

func NewMessageService(store domain.MessageStore, eventStore domain.EventStore) *MessageService {
	return &MessageService{
		store:      store,
		eventStore: eventStore,
	}
}

func (service *MessageService) GetConversation(ctx context.Context, sender, receiver string) (*domain.Conversation, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConversation service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetConversation(ctx, sender, receiver)
}

func (service *MessageService) GetAllConversationsForUser(ctx context.Context, user string) ([]*domain.Conversation, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllConversationsForUser service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetAllConversationsForUser(ctx, user)
}

func (service *MessageService) NewMessage(ctx context.Context, message *domain.Message, sender string) (*domain.Conversation, error) {
	span := tracer.StartSpanFromContext(ctx, "NewMessage service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.NewMessage(ctx, message, sender)
}

func (service *MessageService) GetConversationById(ctx context.Context, id primitive.ObjectID) (*domain.Conversation, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConversationById service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetConversationById(ctx, id)
}

func (service *MessageService) NewEvent(event *domain.Event) (*domain.Event, error) {
	_, err := service.eventStore.NewEvent(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (service *MessageService) GetAllEvents() ([]*domain.Event, error) {
	return service.eventStore.GetAllEvents()
}
