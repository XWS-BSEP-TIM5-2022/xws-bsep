package persistence

import (
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"strings"
)

func checkIfUserExist(userID string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (existing_user:USER) WHERE existing_user.userID = $userID RETURN existing_user.userID",
		map[string]interface{}{"userID": userID})

	if result != nil && result.Next() && result.Record().Values[0] == userID {
		return true
	}
	return false
}

func checkIfSkillExist(skillName string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (existing_skill:SKILL) WHERE toUpper(existing_skill.name) = $name RETURN toUpper(existing_skill.name)",
		map[string]interface{}{"name": strings.ToUpper(skillName)})

	if result != nil && result.Next() && result.Record().Values[0] == strings.ToUpper(skillName) {
		return true
	}
	return false
}

func checkIfExperienceExist(expName string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (e:POSITION) WHERE toUpper(e.headLine) = $headline RETURN toUpper(e.headLine)",
		map[string]interface{}{"headline": strings.ToUpper(expName)})

	if result != nil && result.Next() && result.Record().Values[0] == strings.ToUpper(expName) {
		return true
	}
	return false
}

func checkIfRelationshipExist(userID, skillName string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (u1:USER) WHERE u1.userID=$uIDa "+
			"MATCH (s:SKILL) WHERE s.name=$name "+
			"MATCH (u1)-[r:KNOWS]->(s) "+
			"RETURN r.date ",
		map[string]interface{}{"uIDa": userID, "name": skillName})

	if result != nil && result.Next() {
		return true
	}
	return false
}

func checkIfExpRelationshipExist(userID, expName string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (u1:USER) WHERE u1.userID=$uIDa "+
			"MATCH (s:POSITION) WHERE s.headLine=$headLine "+
			"MATCH (u1)-[r:WORKED]->(s) "+
			"RETURN s.headLine",
		map[string]interface{}{"uIDa": userID, "headLine": expName})

	if result != nil && result.Next() {
		return true
	}
	return false
}

func jobOfferExist(jobId string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (j:JOB) WHERE j.jobID = $id RETURN j.jobID",
		map[string]interface{}{"id": jobId})

	if result != nil && result.Next() && result.Record().Values[0] == jobId {
		return true
	}
	return false
}

func checkIfJobRelationshipExist(jobID, skillName string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (j:JOB) WHERE j.jobID=$jobID "+
			"MATCH (s:SKILL) WHERE s.name=$name "+
			"MATCH (j)-[r:NEEDS]->(s) "+
			"RETURN r ",
		map[string]interface{}{"jobID": jobID, "name": skillName})

	if result != nil && result.Next() {
		return true
	}
	return false
}

func checkIfJobPositionRelationshipExist(jobID, position string, transaction neo4j.Transaction) bool {
	result, _ := transaction.Run(
		"MATCH (j:JOB) WHERE j.jobID=$jobID "+
			"MATCH (s:POSITION) WHERE s.headLine=$position "+
			"MATCH (j)-[r:INCLUDES]->(s) "+
			"RETURN r ",
		map[string]interface{}{"jobID": jobID, "position": position})

	if result != nil && result.Next() {
		return true
	}
	return false
}

func getJobRecommendations(userID string, transaction neo4j.Transaction) ([]*domain.PostsID, error) {
	result, err := transaction.Run(
		"MATCH (u1:USER)-[:WORKED]->(u2:POSITION)<-[:INCLUDES]-(u3:JOB) "+
			"WHERE u1.userID=$uID "+
			"RETURN distinct u3.jobID "+
			"LIMIT 20 ",
		map[string]interface{}{"uID": userID})

	if err != nil {
		return nil, err
	}

	fmt.Println(result)

	var recommendation []*domain.PostsID
	for result.Next() {
		recommendation = append(recommendation, &domain.PostsID{Id: result.Record().Values[0].(string)})
	}
	return recommendation, nil
}
