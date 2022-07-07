package persistence

import (
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
