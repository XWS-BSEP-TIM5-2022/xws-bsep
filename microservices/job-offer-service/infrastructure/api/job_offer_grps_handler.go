package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/job_offer_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/application"
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
