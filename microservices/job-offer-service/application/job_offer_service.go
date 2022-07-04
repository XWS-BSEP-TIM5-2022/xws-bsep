package application

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/domain"
)

type JobOfferService struct {
	store domain.JobOfferStore
}

func NewJobOfferService(store domain.JobOfferStore) *JobOfferService {
	return &JobOfferService{
		store: store,
	}
}
