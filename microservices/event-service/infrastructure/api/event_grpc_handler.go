package api

import (
	"context"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/event_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/event_service/application"
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
func (handler *EventHandler) GetAllEvents(ctx context.Context, request *pb.GetAllEventsRequest) (*pb.GetAllEventsResponse, error) {

	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)

	events, err := handler.service.GetAllEvents()

	handler.CustomLogger.InfoLogger.Info("Get all events for admin with ID: " + userId)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error while getting events for admin: " + userId)
		return nil, err
	}

	var finalEvents []*pb.Event

	for _, event := range events {
		finalEvents = append(finalEvents, mapEvent(event))
	}

	response := &pb.GetAllEventsResponse{
		Events: finalEvents,
	}

	handler.CustomLogger.SuccessLogger.Info("Get all events for admin with ID: " + userId + " successfully done")
	return response, nil

}
