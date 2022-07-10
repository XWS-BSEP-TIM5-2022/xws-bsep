package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/message_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func mapConversation(conversation *domain.Conversation) *pb.Conversation {
	conversationPb := &pb.Conversation{
		Id:    conversation.Id.Hex(),
		User1: conversation.User1,
		User2: conversation.User2,
	}

	for _, message := range conversation.Messages {
		conversationPb.Messages = append(conversationPb.Messages, &pb.Message{
			Id:       message.Id.Hex(),
			Receiver: message.Receiver,
			Content:  message.Content,
			Time:     message.Time.String(),
		})
	}

	return conversationPb
}

func mapEvent(event *domain.Event) *pb.Event {
	eventPb := &pb.Event{
		Id:     event.Id.Hex(),
		UserId: event.UserId,
		Text:   event.Text,
		Date:   event.Date.String(),
	}

	return eventPb
}

func mapMessage(message *domain.Message) *pb.Message {
	messagePb := &pb.Message{
		Id:       message.Id.Hex(),
		Receiver: message.Receiver,
		Content:  message.Content,
		Time:     message.Time.String(),
	}

	return messagePb
}

func mapInsertMessage(message *pb.Message) *domain.Message {
	//receiver, _ := primitive.ObjectIDFromHex(message.Receiver)

	messagePb := &domain.Message{
		Id:       primitive.NewObjectID(),
		Receiver: message.Receiver,
		Content:  removeMalicious(message.Content),
		Time:     time.Now(),
	}

	return messagePb
}

func removeMalicious(value string) string {

	var lenId = len(value)
	var checkId = ""
	for i := 0; i < lenId; i++ {
		char := string(value[i])
		if char != "$" {
			checkId = checkId + char
		}
	}
	return checkId
}
