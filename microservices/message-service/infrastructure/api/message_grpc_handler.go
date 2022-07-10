package api

import (
	"context"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/message_service"
	notification "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"regexp"
	"time"
)

type MessageHandler struct {
	service      *application.MessageService
	CustomLogger *CustomLogger
	pb.UnimplementedMessageServiceServer
	notificationServiceClient notification.NotificationServiceClient
	userServiceClient         user.UserServiceClient
}

func NewMessageHandler(service *application.MessageService, notificationServiceClient notification.NotificationServiceClient,
	userServiceClient user.UserServiceClient) *MessageHandler {
	CustomLogger := NewCustomLogger()
	return &MessageHandler{
		service:                   service,
		CustomLogger:              CustomLogger,
		notificationServiceClient: notificationServiceClient,
		userServiceClient:         userServiceClient,
	}
}

func (handler *MessageHandler) GetConversationById(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConversationById")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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
	conversation, err := handler.service.GetConversationById(ctx, objectId)
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
	span := tracer.StartSpanFromContext(ctx, "GetConversation")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(ctx, span)

	senderId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	conversation, err := handler.service.GetConversation(ctx, senderId, request.Receiver)

	handler.CustomLogger.InfoLogger.Info("Get conversation for user with ID: " + senderId + " and user with ID: " + request.Receiver)
	if err != nil {

		handler.CustomLogger.ErrorLogger.Info("Error getting conversation for user with ID: " + senderId + " and user with ID: " + request.Receiver)
		return nil, err
	}

	response := &pb.GetConversationResponse{
		Conversation: mapConversation(conversation),
	}

	handler.CustomLogger.SuccessLogger.Info("Get conversation for user with ID: " + senderId + " and user with ID: " + request.Receiver + " success!")
	return response, nil
}

func (handler *MessageHandler) GetAllConversationsForUser(ctx context.Context, request *pb.GetAllConversationsForUserRequest) (*pb.GetAllConversationsForUserResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllConversationsForUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(ctx, span)

	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	conversations, err := handler.service.GetAllConversationsForUser(ctx, userId)

	handler.CustomLogger.InfoLogger.Info("Get all conversations for user with ID: " + userId)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error while getting conversations for user: " + userId)
		return nil, err
	}

	var finalConversations []*pb.Conversation

	for _, conversation := range conversations {
		finalConversations = append(finalConversations, mapConversation(conversation))
	}

	response := &pb.GetAllConversationsForUserResponse{
		Conversations: finalConversations,
	}

	handler.CustomLogger.SuccessLogger.Info("Get all conversations for user with ID: " + userId + " successfully done")
	return response, nil
}

func (handler *MessageHandler) NewMessage(ctx context.Context, request *pb.NewMessageRequest) (*pb.NewMessageResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "NewMessage")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(ctx, span)

	sender := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	conversation, err := handler.service.NewMessage(ctx, mapInsertMessage(request.Message), sender)

	handler.CustomLogger.InfoLogger.Info("New message from user with ID: " + sender + " to user with ID: " + request.Message.Receiver)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Info("Error while sending message from user with ID: " + sender + " to user with ID: " + request.Message.Receiver)
		return nil, err
	}

	response := &pb.NewMessageResponse{
		Conversation: mapConversation(conversation),
	}

	successLogText := "New message from user with ID: " + sender + " to user with ID: " + request.Message.Receiver + " sent!"

	handler.CustomLogger.SuccessLogger.Info(successLogText)

	event := domain.Event{
		Id:     primitive.NewObjectID(),
		UserId: sender,
		Text:   successLogText,
		Date:   time.Now(),
	}
	handler.service.NewEvent(&event)

	// slanje notifikacija
	current_user, _ := handler.userServiceClient.Get(ctx, &user.GetRequest{Id: sender})
	reciever, _ := handler.userServiceClient.Get(ctx, &user.GetRequest{Id: request.Message.Receiver})
	if reciever.User.MessageNotification == true {
		notificationRequest := &notification.InsertNotificationRequest{}
		notificationRequest.Notification = &notification.Notification{}
		notificationRequest.Notification.Type = notification.Notification_NotificationTypeEnum(0)
		notificationRequest.Notification.Text = "User " + current_user.User.Name + " " + current_user.User.LastName + " messaged you"
		notificationRequest.Notification.UserId = request.Message.Receiver
		handler.notificationServiceClient.Insert(ctx, notificationRequest)
	}

	return response, nil

}

func (handler *MessageHandler) GetAllEvents(ctx context.Context, request *pb.GetAllEventsRequest) (*pb.GetAllEventsResponse, error) {

	events, err := handler.service.GetAllEvents()

	handler.CustomLogger.InfoLogger.Info("Get all events for admin.")

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error while getting events for admin")
		return nil, err
	}

	var finalEvents []*pb.Event

	for _, event := range events {
		finalEvents = append(finalEvents, mapEvent(event))
	}

	response := &pb.GetAllEventsResponse{
		Events: finalEvents,
	}

	handler.CustomLogger.SuccessLogger.Info("Get all events for admin successfully done")
	return response, nil

}
