package persistence

import (
	"fmt"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/domain"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"time"
)

type ConnectionDBStore struct {
	connectionDB *neo4j.Driver
}

func NewConnectionDBStore(client *neo4j.Driver) domain.ConnectionStore {
	return &ConnectionDBStore{
		connectionDB: client,
	}
}

func (store *ConnectionDBStore) Register(userID string, isPublic bool) (*pb.ActionResult, error) {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		actionResult := &pb.ActionResult{}

		if checkIfUserExist(userID, transaction) {
			actionResult.Status = 406
			actionResult.Msg = "error user with ID:" + userID + " already exist"
			return actionResult, nil
		}

		_, err := transaction.Run(
			"CREATE (new_user:USER{userID:$userID, isPublic:$isPublic})",
			map[string]interface{}{"userID": userID, "isPublic": isPublic})

		if err != nil {
			actionResult.Msg = "error while creating new node with ID:" + userID
			actionResult.Status = 501
			return actionResult, err
		}

		actionResult.Msg = "successfully created new node with ID:" + userID
		actionResult.Status = 201

		return actionResult, err
	})

	return result.(*pb.ActionResult), err
}

func (store *ConnectionDBStore) GetFriends(userID string) ([]domain.UserConn, error) {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	friends, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (this_user:USER) -[:FRIEND]-> (my_friend:USER) WHERE this_user.userID=$uID RETURN my_friend.userID, my_friend.isPublic",
			map[string]interface{}{"uID": userID})

		if err != nil {
			return nil, err
		}

		var friends []domain.UserConn
		for result.Next() {
			friends = append(friends, domain.UserConn{UserID: result.Record().Values[0].(string), IsPublic: result.Record().Values[1].(bool)})
		}
		return friends, nil

	})
	if err != nil {
		return nil, err
	}

	return friends.([]domain.UserConn), nil
}

func (store *ConnectionDBStore) AddConnection(userIDa, userIDb string) (*pb.ActionResult, error) {
	/*
				Dodavanje novog prijatelja je moguce ako:
		         - userA i userB postoji
				 - userA nije prijatelj sa userB
				 - userA nije blokirao userB
			   	 - userA nije blokiran od strane userB
	*/

	if userIDa == userIDb {
		return &pb.ActionResult{Msg: "userIDa is same as userIDb", Status: 400}, nil
	}

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		actionResult := &pb.ActionResult{Msg: "msg", Status: 0}

		//ako ne postoji userA, kreira ga
		if !checkIfUserExist(userIDa, transaction) {
			_, err := transaction.Run(
				"CREATE (new_user:USER{userID:$userID, isPublic:$isPublic})",
				map[string]interface{}{"userID": userIDa, "isPublic": true})

			if err != nil {
				actionResult.Msg = "error while creating new node with ID:" + userIDa
				actionResult.Status = 501
				return actionResult, err
			}
		}
		//ako ne postoji userB, kreira ga
		if !checkIfUserExist(userIDb, transaction) {
			_, err := transaction.Run(
				"CREATE (new_user:USER{userID:$userID, isPublic:$isPublic})",
				map[string]interface{}{"userID": userIDb, "isPublic": false})

			if err != nil {
				actionResult.Msg = "error while creating new node with ID:" + userIDb
				actionResult.Status = 501
				return actionResult, err
			}
		}

		if checkIfUserExist(userIDa, transaction) && checkIfUserExist(userIDb, transaction) {
			if checkIfFriendExist(userIDa, userIDb, transaction) || checkIfFriendExist(userIDb, userIDa, transaction) {
				actionResult.Msg = "users are already friends"
				actionResult.Status = 400 //bad request
				return actionResult, nil
			} else {
				//if checkIfBlockExist(userIDa, userIDb, transaction) || checkIfBlockExist(userIDb, userIDa, transaction) {
				//	actionResult.Msg = "block already exist"
				//	actionResult.Status = 400 //bad request
				//	return actionResult, nil
				//} else {

				//ako je userB public, odmah ce kreirati konekciju
				if !checkIfPublicUser(userIDb, transaction) {
					dateNow := time.Now().Local().Unix()
					result, err := transaction.Run(
						"MATCH (u1:USER) WHERE u1.userID=$uIDa "+
							"MATCH (u2:USER) WHERE u2.userID=$uIDb "+
							"CREATE (u1)-[r1:FRIEND {date: $dateNow}]->(u2) "+
							"CREATE (u2)-[r2:FRIEND {date: $dateNow}]->(u1) "+
							"RETURN r1.date, r2.date",
						map[string]interface{}{"uIDa": userIDa, "uIDb": userIDb, "dateNow": dateNow})

					if err != nil || result == nil {
						actionResult.Msg = "error while creating new friends IDa:" + userIDa + " and IDb:" + userIDb
						actionResult.Status = 501
						return actionResult, err
					}
				} else {
					fmt.Println("USER JE PRIVATE!")
					actionResult.Msg = "user with id: " + userIDb + "is private!"
					actionResult.Status = 501
				}
			}
			//}
		} else {
			actionResult.Msg = "user does not exist"
			actionResult.Status = 400 //bad request
			return actionResult, nil
		}

		actionResult.Msg = "successfully created new friends IDa:" + userIDa + " and IDb:" + userIDb
		actionResult.Status = 201

		return actionResult, nil
	})

	if result == nil {
		return &pb.ActionResult{Msg: "error", Status: 500}, err
	} else {
		return result.(*pb.ActionResult), err
	}
}
