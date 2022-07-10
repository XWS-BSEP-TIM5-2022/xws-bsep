package application

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/domain"
)

type JobOfferService struct {
	store      domain.JobOfferStore
	eventStore domain.EventStore
}

func NewJobOfferService(store domain.JobOfferStore, eventStore domain.EventStore) *JobOfferService {
	return &JobOfferService{
		store:      store,
		eventStore: eventStore,
	}
}

func (service *JobOfferService) GetRecommendations(user *domain.User, jobOffers []*domain.Post) ([]*domain.PostsID, error) {

	var recommendations []*domain.PostsID

	recommendations, err := service.store.GetRecommendations(user, jobOffers)
	if err != nil {
		return nil, nil
	}
	//for _, r := range recommendations {
	//	recommendations = append(recommendations, r)
	//}
	return recommendations, nil
}

func (service *JobOfferService) NewEvent(event *domain.Event) (*domain.Event, error) {
	_, err := service.eventStore.NewEvent(event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (service *JobOfferService) GetAllEvents() ([]*domain.Event, error) {
	return service.eventStore.GetAllEvents()
}
