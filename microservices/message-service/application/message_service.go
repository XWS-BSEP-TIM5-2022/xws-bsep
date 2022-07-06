package application

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageService struct {
	store domain.MessageStore
}

func NewMessageService(store domain.MessageStore) *MessageService {
	return &MessageService{
		store: store,
	}
}

func (service *MessageService) GetConversation(sender, receiver string) (*domain.Conversation, error) {
	return service.store.GetConversation(sender, receiver)
}

func (service *MessageService) GetAllConversationsForUser(user string) ([]*domain.Conversation, error) {
	return service.store.GetAllConversationsForUser(user)
}

func (service *MessageService) NewMessage(message *domain.Message, sender string) (*domain.Conversation, error) {
	return service.store.NewMessage(message, sender)
}

func (service *MessageService) GetConversationById(id primitive.ObjectID) (*domain.Conversation, error) {
	return service.store.GetConversationById(id)
}
