package application

import (
	"context"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
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

func (service *JobOfferService) GetRecommendations(ctx context.Context, user *domain.User, jobOffers []*domain.Post) ([]*domain.PostsID, error) {
	span := tracer.StartSpanFromContext(ctx, "GetRecommendations service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	var recommendations []*domain.PostsID

	recommendations, err := service.store.GetRecommendations(ctx, user, jobOffers)
	if err != nil {
		return nil, nil
	}
	//for _, r := range recommendations {
	//	recommendations = append(recommendations, r)
	//}
	return recommendations, nil
}
