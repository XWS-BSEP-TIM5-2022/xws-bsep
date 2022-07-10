package api

import (
	"context"
	"fmt"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/job_offer_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/domain"
	"strconv"
)

type JobOfferHandler struct {
	pb.UnimplementedJobOfferServiceServer
	service      *application.JobOfferService
	CustomLogger *CustomLogger
}

func NewJobOfferHandler(service *application.JobOfferService) *JobOfferHandler {
	CustomLogger := NewCustomLogger()
	return &JobOfferHandler{
		service:      service,
		CustomLogger: CustomLogger,
	}
}

func (handler *JobOfferHandler) GetRecommendations(ctx context.Context, request *pb.GetRequest) (*pb.Recommendations, error) {
	handler.CustomLogger.ErrorLogger.Info("usao sam")

	var jobs []*domain.Post

	for _, job := range request.DTO.JobOffers.GetJobOffers() {
		domainJob := mapJobOffer(job)
		jobs = append(jobs, domainJob)
	}
	user := mapUser(request.GetDTO().User)

	recommendations, err := handler.service.GetRecommendations(user, jobs)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Job recommendations for user with ID: " + request.DTO.User.Id + " not found")
		return nil, err
	}
	response := &pb.Recommendations{}
	for _, rec := range recommendations {
		response.JobOffers = append(response.JobOffers, mapRecommendations(rec))

		fmt.Println(rec)
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(recommendations)) + " recommendations for user with ID: " + request.DTO.User.Id)
	return response, nil
}

func (handler *JobOfferHandler) GetAllEvents(ctx context.Context, request *pb.GetAllEventsRequest) (*pb.GetAllEventsResponse, error) {

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
