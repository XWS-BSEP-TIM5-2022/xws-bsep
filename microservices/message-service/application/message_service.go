package application

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/domain"
)

type MessageService struct {
	store domain.MessageStore
}

func NewMessageService(store domain.MessageStore) *MessageService {
	return &MessageService{
		store: store,
	}
}
