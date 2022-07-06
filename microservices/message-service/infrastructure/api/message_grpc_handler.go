package api

import (
	"context"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/message_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"regexp"
)

type MessageHandler struct {
	service      *application.MessageService
	CustomLogger *CustomLogger
	pb.UnimplementedMessageServiceServer
}

func NewMessageHandler(service *application.MessageService) *MessageHandler {
	CustomLogger := NewCustomLogger()
	return &MessageHandler{
		service:      service,
		CustomLogger: CustomLogger,
	}
}

func (handler *MessageHandler) GetConversationById(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := removeMalicious(request.Id)
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	requestId := re.ReplaceAllString(request.Id, " ")

	handler.CustomLogger.InfoLogger.WithField("id", requestId).Info("Getting conversation by id: " + requestId)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	conversation, err := handler.service.GetConversationById(objectId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Conversation with ID:" + objectId.Hex() + " not found")
		return nil, err
	}
	conversationPb := mapConversation(conversation)
	response := &pb.GetResponse{
		Conversation: conversationPb,
	}
	handler.CustomLogger.SuccessLogger.Info("Conversation by ID:" + objectId.Hex() + " received successfully")
	return response, nil
}

func (handler *MessageHandler) GetConversation(ctx context.Context, request *pb.GetConversationRequest) (*pb.GetConversationResponse, error) {

	senderId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	conversation, err := handler.service.GetConversation(senderId, request.Receiver)

	if err != nil {
		return nil, err
	}

	response := &pb.GetConversationResponse{
		Conversation: mapConversation(conversation),
	}
	return response, nil
}

func (handler *MessageHandler) GetAllConversationsForUser(ctx context.Context, request *pb.GetAllConversationsForUserRequest) (*pb.GetAllConversationsForUserResponse, error) {

	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	conversations, err := handler.service.GetAllConversationsForUser(userId)

	if err != nil {
		return nil, err
	}

	var finalConversations []*pb.Conversation

	for _, messHistory := range conversations {
		finalConversations = append(finalConversations, mapConversation(messHistory))
	}

	response := &pb.GetAllConversationsForUserResponse{
		Conversations: finalConversations,
	}

	return response, nil
}

func (handler *MessageHandler) NewMessage(ctx context.Context, request *pb.NewMessageRequest) (*pb.NewMessageResponse, error) {

	sender := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	conversation, err := handler.service.NewMessage(mapInsertMessage(request.Message), sender)

	if err != nil {
		return nil, err
	}

	response := &pb.NewMessageResponse{
		Conversation: mapConversation(conversation),
	}

	return response, nil

}
