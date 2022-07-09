package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/event_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/event_service/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"regexp"
)

type EventHandler struct {
	service      *application.EventService
	CustomLogger *CustomLogger
	pb.UnimplementedEventServiceServer
}

func NewEventHandler(service *application.EventService) *EventHandler {
	CustomLogger := NewCustomLogger()
	return &EventHandler{
		service:      service,
		CustomLogger: CustomLogger,
	}
}

func (handler *EventHandler) GetById(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := removeMalicious(request.Id)
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	requestId := re.ReplaceAllString(request.Id, " ")

	handler.CustomLogger.InfoLogger.WithField("id", requestId).Info("Getting event by id: " + requestId)

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	event, err := handler.service.GetById(objectId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Event with ID:" + objectId.Hex() + " not found")
		return nil, err
	}
	eventPb := mapEvent(event)
	response := &pb.GetResponse{
		Event: eventPb,
	}
	handler.CustomLogger.SuccessLogger.Info("Event by ID:" + objectId.Hex() + " received successfully")
	return response, nil
}
