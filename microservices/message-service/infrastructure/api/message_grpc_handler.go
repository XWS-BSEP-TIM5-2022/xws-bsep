package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/message_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/application"
)

type MessageHandler struct {
	service *application.MessageService
	pb.UnimplementedMessageServiceServer
}

func NewMessageHandler(service *application.MessageService) *MessageHandler {
	return &MessageHandler{
		service: service,
	}
}

