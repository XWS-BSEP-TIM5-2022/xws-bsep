package api

import (
	"context"
	"encoding/json"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/infrastructure/services"
	authGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	connectionGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	offerGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/job_offer_service"
	messageGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/message_service"
	notificationGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	postGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	userGw "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
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

	//endpointName := "GetAllEvents"
	//span := tracer.StartSpanFromContext(r.Context(), endpointName)
	//defer span.Finish()
	//ctx := tracer.ContextWithSpan(r.Context(), span)

	messageClient := services.NewMessageClient(handler.messageClientAddress)
	connectionClient := services.NewConnectionClient(handler.connectionClientAddress)
	userClient := services.NewUserClient(handler.userClientAddress)
	authClien := services.NewAuthClient(handler.authClientAddress)
	postClient := services.NewPostClient(handler.postClientAddress)
	offerClient := services.NewJobOfferClient(handler.jobOfferClientAddress)
	notificationClient := services.NewNotificationClient(handler.notificationClientAddress)

	finalEvents, err := messageClient.GetAllEvents(context.TODO(), &messageGw.GetAllEventsRequest{})

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error getting all events for message service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	connectionEvents, err := connectionClient.GetAllEvents(context.TODO(), &connectionGw.GetAllEventsRequest{})

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

	userEvents, err := userClient.GetAllEvents(context.TODO(), &userGw.GetAllEventsRequest{})

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

	postEvents, err := postClient.GetAllEvents(context.TODO(), &postGw.GetAllEventsRequest{})

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error getting all events for post service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, event := range postEvents.Events {
		finalEvents.Events = append(finalEvents.Events, &messageGw.Event{
			Id:     event.Id,
			UserId: event.UserId,
			Text:   event.Text,
			Date:   event.Date,
		})
	}

	authEvents, err := authClien.GetAllEvents(context.TODO(), &authGw.GetAllEventsRequest{})

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error getting all events for auth service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, event := range authEvents.Events {
		finalEvents.Events = append(finalEvents.Events, &messageGw.Event{
			Id:     event.Id,
			UserId: event.UserId,
			Text:   event.Text,
			Date:   event.Date,
		})
	}

	offerEvents, err := offerClient.GetAllEvents(context.TODO(), &offerGw.GetAllEventsRequest{})

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error getting all events for job offer service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, event := range offerEvents.Events {
		finalEvents.Events = append(finalEvents.Events, &messageGw.Event{
			Id:     event.Id,
			UserId: event.UserId,
			Text:   event.Text,
			Date:   event.Date,
		})
	}

	notificationEvents, err := notificationClient.GetAllEvents(context.TODO(), &notificationGw.GetAllEventsRequest{})

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Error getting all events for notification service")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, event := range notificationEvents.Events {
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
