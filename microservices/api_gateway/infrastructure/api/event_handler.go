package api

import (
	"context"
	"encoding/json"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/services"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"

	connectionGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
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
	err := mux.HandlePath("GET", "/events", handler.GetAllEvents)
	if err != nil {
		panic(err)
	}
}

func (handler *EventHandler) GetAllEvents(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {

	span := tracer.StartSpanFromContext(nil, "GetAllEvents") //TODO: STA S OVIM?
	defer span.Finish()

	ctx := tracer.ContextWithSpan(context.Background(), span)

	messageClient := services.NewMessageClient(handler.messageClientAddress)
	connectionClient := services.NewConnectionClient(handler.connectionClientAddress)

	var finalEvents []*domain.Event

	messageEvents, err := messageClient.GetAllEvents(ctx, &messageGw.GetAllEventsRequest{})

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error getting all events for message service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, event := range messageEvents.Events {
		finalEvents = append(finalEvents, &domain.Event{
			Id:     event.Id,
			UserId: event.UserId,
			Text:   event.Text,
			Date:   event.Date,
		})
	}

	connectionEvents, err := connectionClient.GetAllEvents(ctx, &connectionGw.GetAllEventsRequest{})

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error getting all events for connection service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, event := range connectionEvents.Events {
		finalEvents = append(finalEvents, &domain.Event{
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
