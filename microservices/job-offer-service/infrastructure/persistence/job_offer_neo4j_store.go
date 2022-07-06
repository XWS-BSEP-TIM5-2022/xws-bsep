package persistence

import (
	"fmt"
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

func (store *JobOfferDBStore) GetRecommendations(user *domain.User, jobOffers []*domain.Post) ([]*domain.Post, error) {
	fmt.Println(user)

	session := (*store.jobOfferDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		//ako ne postoji korisnik, dodaje ga
		if !checkIfUserExist(user.Id.String(), transaction) {
			_, err := transaction.Run(
				"CREATE (new_user:USER{userID:$userID, isPublic:$isPublic, name : $Name, lastName: $LastName, email: $Email})",
				map[string]interface{}{"userID": user.Id.String(), "Name": user.Name, "LastName": user.LastName, "Email": user.Email, "isPublic": user.IsPublic})

			if err != nil {
				return nil, err
			}

		}

		//ako ne postoje vjestine, dodaje ih
		for _, s := range user.Skills {
			if !checkIfSkillExist(s.Name, transaction) {
				_, err := transaction.Run(
					"CREATE (new_skill:SKILL{name : $Name})",
					map[string]interface{}{"Name": s.Name})

				if err != nil {
					return nil, err
				}

			}

			//ako korisnik nije povezan sa vjestinama, dodaje ih
			if !checkIfRelationshipExist(user.Id.String(), s.Name, transaction) {
				result, err := transaction.Run(
					"MATCH (u:USER) WHERE u.userID=$uIDa "+
						"MATCH (s:SKILL) WHERE s.name=$name "+
						"CREATE (u)-[r:KNOWS]->(s) "+
						"RETURN u.userID",
					map[string]interface{}{"uIDa": user.Id.String(), "name": s.Name})
				if err != nil {
					return nil, err
				}

				fmt.Println(result)
			}
		}

		return nil, nil

	})

	fmt.Println(result)
	fmt.Println(err)

	return nil, nil
}
