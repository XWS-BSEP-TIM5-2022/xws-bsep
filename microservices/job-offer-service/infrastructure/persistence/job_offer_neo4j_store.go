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

		//TODO: sredi povratne vrijednosti
		//TODO: dodaj sva polja kod korisnika
		//TODO: dodaj jobOffere
		//TODO: spoj jobOffere sa skilovima
		//TODO: upit za preporuke
		//TODO: sve isto za radno iskustvo!!!!!!

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

		//ako ne postoji iskustvo, dodaje ga
		for _, s := range user.Experience {
			if !checkIfExperienceExist(s.Headline, transaction) {
				_, err := transaction.Run(
					"CREATE (new_exp:POSITION{headLine : $headLine}) ",
					map[string]interface{}{"headLine": s.Headline})

				if err != nil {
					return nil, err
				}

			}

			//ako korisnik nije povezan sa vjestinama, dodaje ih
			if !checkIfExpRelationshipExist(user.Id.String(), s.Headline, transaction) {
				result, err := transaction.Run(
					"MATCH (u:USER) WHERE u.userID=$uIDa "+
						"MATCH (s:POSITION) WHERE s.headLine=$headLine "+
						"CREATE (u)-[r:WORKED]->(s) "+
						"RETURN u.userID",
					map[string]interface{}{"uIDa": user.Id.String(), "headLine": s.Headline})
				if err != nil {
					return nil, err
				}

				fmt.Println(result)
			}
		}

		//ako ne postoji job offer, dodaje ga
		//
		for _, job := range jobOffers {

			if !jobOfferExist(job.Id.String(), transaction) {
				_, err := transaction.Run(
					"CREATE (new_job:JOB{position:$position, jobID:$jobID, text:$Text, preconditions: $preconditions})",
					map[string]interface{}{"jobID": job.Id.String(), "Text": job.Text, "preconditions": job.JobOffer.Preconditions, "position": job.JobOffer.Position.Name})

				if err != nil {
					return nil, err
				}
			}

			//ako jobOffer nije povezan sa vjestinama, povezuje ih
			if !checkIfJobRelationshipExist(job.Id.String(), job.JobOffer.Preconditions, transaction) {
				result, err := transaction.Run(
					"MATCH (j:JOB) WHERE j.jobID=$jobID "+
						"MATCH (s:SKILL) WHERE s.name=$name "+
						"CREATE (j)-[r:NEEDS]->(s) "+
						"RETURN j.jobID",
					map[string]interface{}{"jobID": job.Id.String(), "name": job.JobOffer.Preconditions})
				if err != nil {
					return nil, err
				}

				fmt.Println(result)
			}

			//ako jobOffer nije povezan sa pozicijom, povezuje ih
			if !checkIfJobPositionRelationshipExist(job.Id.String(), job.JobOffer.Position.Name, transaction) {
				result, err := transaction.Run(
					"MATCH (j:JOB) WHERE j.jobID=$jobID "+
						"MATCH (s:POSITION) WHERE s.headLine=$position "+
						"CREATE (j)-[r:INCLUDES]->(s) "+
						"RETURN j.jobID",
					map[string]interface{}{"jobID": job.Id.String(), "position": job.JobOffer.Position.Name})
				if err != nil {
					return nil, err
				}

				fmt.Println(result)
			}
		}

		//var recommendation []*domain.PostsID

		//jobsRecommendations, err1 := getJobRecommendations(user.Id.String(), transaction)
		//if err1 != nil {
		//	return recommendation, err1
		//}
		//
		//for _, recommend := range jobsRecommendations {
		//	recommendation = append(recommendation, recommend)
		//}
		//
		//return recommendation, err1
		return nil, nil
	})

	fmt.Println(result)
	fmt.Println(err)
	//if err != nil || result == nil {
	//	return nil, err
	//}

	//return result.([]*domain.PostsID), nil
	return nil, nil
}
