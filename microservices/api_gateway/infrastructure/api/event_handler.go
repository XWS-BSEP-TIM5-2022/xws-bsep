package api

import (
	"encoding/json"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/services"
	connectionGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	userGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"

	messageGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/message_service"
)

type EventHandler struct {
	authClientAddress         string
	userClientAddress         string
	postClientAddress         string
	connectionClientAddress   string
	messageClientAddress      string
	notificationClientAddress string
	jobOfferClientAddress     string
	CustomLogger              *CustomLogger
	//TODO JH: tracing
}

func NewEventHandler(authClientAddress string, userClientAddress string, postClientAddress string, connectionClientAddress string, messageClientAddress string, notificationClientAddress string, jobOfferClientAddress string) Handler {
	CustomLogger := NewCustomLogger()
	return &EventHandler{
		authClientAddress:         authClientAddress,
		userClientAddress:         userClientAddress,
		postClientAddress:         postClientAddress,
		connectionClientAddress:   connectionClientAddress,
		messageClientAddress:      messageClientAddress,
		notificationClientAddress: notificationClientAddress,
		jobOfferClientAddress:     jobOfferClientAddress,
		CustomLogger:              CustomLogger,
	}
}

func (handler *EventHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/GetAllEvents", handler.GetAllEvents)
	if err != nil {
		panic(err)
	}
}

func (handler *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

	endpointName := "GetAllEvents"
	span := tracer.StartSpanFromContext(r.Context(), endpointName)
	defer span.Finish()
	ctx := tracer.ContextWithSpan(r.Context(), span)

	messageClient := services.NewMessageClient(handler.messageClientAddress)
	connectionClient := services.NewConnectionClient(handler.connectionClientAddress)
	userClient := services.NewUserClient(handler.userClientAddress)

	finalEvents, err := messageClient.GetAllEvents(ctx, &messageGw.GetAllEventsRequest{})

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error getting all events for message service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	connectionEvents, err := connectionClient.GetAllEvents(ctx, &connectionGw.GetAllEventsRequest{})

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error getting all events for connection service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, event := range connectionEvents.Events {
		finalEvents.Events = append(finalEvents.Events, &messageGw.Event{
			Id:     event.Id,
			UserId: event.UserId,
			Text:   event.Text,
			Date:   event.Date,
		})
	}

	userEvents, err := userClient.GetAllEvents(ctx, &userGw.GetAllEventsRequest{})

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error getting all events for user service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, event := range userEvents.Events {
		finalEvents.Events = append(finalEvents.Events, &messageGw.Event{
			Id:     event.Id,
			UserId: event.UserId,
			Text:   event.Text,
			Date:   event.Date,
		})
	}

	response, err := json.Marshal(finalEvents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
