package persistence

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type JobOfferDBStore struct {
	jobOfferDB *neo4j.Driver
}

func NewJobOfferDBStore(client *neo4j.Driver) domain.JobOfferStore {
	return &JobOfferDBStore{
		jobOfferDB: client,
	}
}
