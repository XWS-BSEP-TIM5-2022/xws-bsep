package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/notification_service/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"regexp"
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
