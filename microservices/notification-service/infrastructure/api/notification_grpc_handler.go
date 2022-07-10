package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/notification_service/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"regexp"
	"strconv"
)

type NotificationHandler struct {
	service      *application.NotificationService
	CustomLogger *CustomLogger
	pb.UnimplementedNotificationServiceServer
}

func NewNotificationHandler(service *application.NotificationService) *NotificationHandler {
	CustomLogger := NewCustomLogger()
	return &NotificationHandler{
		service:      service,
		CustomLogger: CustomLogger,
	}
}

func (handler *NotificationHandler) GetById(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {

	id := removeMalicious(request.Id)
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	requestId := re.ReplaceAllString(request.Id, " ")

	handler.CustomLogger.InfoLogger.WithField("id", requestId).Info("Getting notification by id: " + requestId)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	notification, err := handler.service.GetById(objectId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Notification with ID:" + objectId.Hex() + " not found")
		return nil, err
	}
	notificationPb := mapNotification(notification)
	response := &pb.GetResponse{
		Notification: notificationPb,
	}
	handler.CustomLogger.SuccessLogger.Info("Notification by ID:" + objectId.Hex() + " received successfully")
	return response, nil
}

func (handler *NotificationHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	posts, err := handler.service.GetAll()
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Get all notifications unsuccessful")
		return nil, err
	}
	response := &pb.GetAllResponse{
		Notifications: []*pb.Notification{},
	}
	for _, post := range posts {
		current := mapNotification(post)
		response.Notifications = append(response.Notifications, current)
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(posts)) + " notification")
	return response, nil
}

func (handler *NotificationHandler) Insert(ctx context.Context, request *pb.InsertNotificationRequest) (*pb.InsertNotificationResponse, error) {
	notification, err := mapInsertNotification(request.Notification)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Notification was not mapped")
		return nil, err
	}

	notification.Id = primitive.NewObjectID()
	success, err := handler.service.Insert(notification)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Notification was not inserted")
		return nil, err
	}
	response := &pb.InsertNotificationResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Notification with ID: " + notification.Id.Hex() + " created")
	return response, err
}

func (handler *NotificationHandler) GetNotificationsByUserId(ctx context.Context, request *pb.GetRequest) (*pb.GetAllResponse, error) {

	id := removeMalicious(request.Id)
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	requestId := re.ReplaceAllString(request.Id, " ")
	handler.CustomLogger.InfoLogger.WithField("id", requestId).Info("Getting notification by id: " + requestId)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	notifications, err := handler.service.GetAllByUser(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Notification with ID:" + objectId.Hex() + " not found")
		return nil, err
	}

	response := &pb.GetAllResponse{
		Notifications: []*pb.Notification{},
	}
	for _, not := range notifications {
		current := mapNotification(not)
		response.Notifications = append(response.Notifications, current)
	}

	handler.CustomLogger.SuccessLogger.Info("Notification by user received successfully")
	return response, nil
}

func (handler *NotificationHandler) GetAllEvents(ctx context.Context, request *pb.GetAllEventsRequest) (*pb.GetAllEventsResponse, error) {

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
