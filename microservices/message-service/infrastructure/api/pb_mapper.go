package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/message_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func mapConversation(conversation *domain.Conversation) *pb.Conversation {
	conversationPb := &pb.Conversation{
		Id:    conversation.Id.Hex(),
		User1: conversation.User1.Hex(),
		User2: conversation.User2.Hex(),
	}

	for _, message := range conversation.Messages {
		conversationPb.Messages = append(conversationPb.Messages, &pb.Message{
			Id:       message.Id.Hex(),
			Receiver: message.Receiver.Hex(),
			Content:  message.Content,
			Time:     timestamppb.New(message.Time),
		})
	}

	return conversationPb
}

func mapInsertConversation(conversation *pb.Conversation) *domain.Conversation {
	id, _ := primitive.ObjectIDFromHex(conversation.Id)
	user1, _ := primitive.ObjectIDFromHex(conversation.User1)
	user2, _ := primitive.ObjectIDFromHex(conversation.User2)

	conversationPb := &domain.Conversation{
		Id:    id,
		User1: user1,
		User2: user2,
	}

	for _, message := range conversation.Messages {

		messId := primitive.NewObjectID()
		receiver, _ := primitive.ObjectIDFromHex(message.Receiver)

		conversationPb.Messages = append(conversationPb.Messages, domain.Message{
			Id:       messId,
			Receiver: receiver,
			Content:  removeMalicious(message.Content),
			Time:     message.Time.AsTime(),
		})
	}

	return conversationPb
}

func mapMessage(message *domain.Message) *pb.Message {
	messagePb := &pb.Message{
		Id:       message.Id.Hex(),
		Receiver: message.Receiver.Hex(),
		Content:  message.Content,
		Time:     timestamppb.New(message.Time),
	}

	return messagePb
}

func mapInsertMessage(message *pb.Message) *domain.Message {
	id, _ := primitive.ObjectIDFromHex(message.Id)
	receiver, _ := primitive.ObjectIDFromHex(message.Receiver)

	messagePb := &domain.Message{
		Id:       id,
		Receiver: receiver,
		Content:  removeMalicious(message.Content),
		Time:     message.Time.AsTime(),
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
