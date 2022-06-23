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
			actionResult.Msg = "User with ID:" + userID + " already exist"
			return actionResult, nil
		}

		_, err := transaction.Run(
			"CREATE (new_user:USER{userID:$userID, isPublic:$isPublic})",
			map[string]interface{}{"userID": userID, "isPublic": isPublic})

		if err != nil {
			actionResult.Msg = "Error while creating new user node with ID:" + userID
			return actionResult, err
		}

		actionResult.Msg = "Successfully created new user node with ID:" + userID

		return actionResult, err
	})

	return result.(*pb.ActionResult), err
}

func (store *ConnectionDBStore) GetConnections(userID string) ([]domain.UserConn, error) {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	friends, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		//result, err := transaction.Run(
		//	"MATCH (this_user:USER) -[:FRIEND]-> (my_friend:USER) WHERE this_user.userID=$uID RETURN my_friend.userID, my_friend.isPublic",
		//	map[string]interface{}{"uID": userID})
		result, err := transaction.Run(
			"MATCH (this_user:USER) -[:FRIEND]-> (my_friend:USER)"+
				" WHERE this_user.userID=$uID "+
				"MATCH (this_user:USER) <-[:FRIEND]- (my_friend:USER)"+
				" WHERE this_user.userID=$uID "+
				"RETURN my_friend.userID, my_friend.isPublic",
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

func (store *ConnectionDBStore) GetRequests(userID string) ([]domain.UserConn, error) {

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	friends, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (u1:USER) WHERE u1.userID=$uID"+
				" MATCH (u2:USER)"+
				"WHERE NOT (u1)-[:FRIEND]->(u2) AND (u2)-[:FRIEND]->(u1)"+
				"RETURN u2.userID, u2.isPublic",
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

func (store *ConnectionDBStore) AddConnection(userIDa string, userIDb string, isPublic bool, isPublicLogged bool) (*pb.AddConnectionResult, error) {
	fmt.Println("Adding new connection")
	fmt.Println(userIDa)
	fmt.Println(userIDb)
	fmt.Println(isPublic)
	fmt.Println(isPublicLogged)

	if userIDa == userIDb {
		return &pb.AddConnectionResult{Msg: "userIDa is same as userIDb", Connected: false, Error: false}, nil
	}

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		actionResult := &pb.AddConnectionResult{Msg: "msg", Connected: false, Error: false}

		//ako ne postoji userA, kreira ga
		if !checkIfUserExist(userIDa, transaction) {
			_, err := transaction.Run(
				"CREATE (new_user:USER{userID:$userID, isPublic:$isPublic})",
				map[string]interface{}{"userID": userIDa, "isPublic": isPublicLogged})

			if err != nil {
				actionResult.Msg = "Error while creating new user node with ID:" + userIDa
				actionResult.Connected = false
				actionResult.Error = true
				return actionResult, err
			}
		}
		//ako ne postoji userB, kreira ga
		if !checkIfUserExist(userIDb, transaction) {
			_, err := transaction.Run(
				"CREATE (new_user:USER{userID:$userID, isPublic:$isPublic})",
				map[string]interface{}{"userID": userIDb, "isPublic": isPublic})

			if err != nil {
				actionResult.Msg = "Error while creating new user node with ID:" + userIDb
				actionResult.Connected = false
				actionResult.Error = true
				return actionResult, err
			}
		}

		if checkIfUserExist(userIDa, transaction) && checkIfUserExist(userIDb, transaction) {
			if checkIfBlockExist(userIDa, userIDb, transaction) || checkIfBlockExist(userIDb, userIDa, transaction) {
				actionResult.Msg = "Users are already blocked"
				actionResult.Connected = true //TODO:provjeri ovo
				actionResult.Error = false
				return actionResult, nil
			}

			if checkIfFriendExist(userIDa, userIDb, transaction) || checkIfFriendExist(userIDb, userIDa, transaction) {
				actionResult.Msg = "Users are already connected"
				actionResult.Connected = true
				actionResult.Error = false
				return actionResult, nil
			} else {

				//ako je userB public, odmah ce kreirati konekciju
				if checkIfPublicUser(userIDb, transaction) {
					dateNow := time.Now().Local().Unix()

					result, err := transaction.Run(
						"MATCH (u1:USER) WHERE u1.userID=$uIDa "+
							"MATCH (u2:USER) WHERE u2.userID=$uIDb "+
							"CREATE (u1)-[r1:FRIEND {date: $dateNow, isApproved: $isApproved}]->(u2) "+
							"CREATE (u2)-[r2:FRIEND {date: $dateNow, isApproved: $isApproved}]->(u1) "+
							"RETURN r1.date, r2.date",
						map[string]interface{}{"uIDa": userIDa, "uIDb": userIDb, "dateNow": dateNow, "isApproved": true})

					actionResult.Msg = "Successfully created new connection IDa:" + userIDa + " and IDb:" + userIDb
					actionResult.Connected = true
					actionResult.Error = false

					if err != nil || result == nil {
						actionResult.Msg = "Error while creating new connection IDa:" + userIDa + " and IDb:" + userIDb
						actionResult.Connected = false
						actionResult.Error = true
						return actionResult, err
					}
				} else {
					//ako je user private kreirace konekciju koja nije odobrena
					dateNow := time.Now().Local().Unix()
					result, err := transaction.Run(
						"MATCH (u1:USER) WHERE u1.userID=$uIDa "+
							"MATCH (u2:USER) WHERE u2.userID=$uIDb "+
							"CREATE (u1)-[r1:FRIEND {date: $dateNow, isApproved: $isApproved}]->(u2) "+
							"RETURN r1.date",
						map[string]interface{}{"uIDa": userIDa, "uIDb": userIDb, "dateNow": dateNow, "isApproved": false})

					actionResult.Msg = "Successfully created new request for connection IDa:" + userIDa + " and IDb:" + userIDb
					actionResult.Connected = false
					actionResult.Error = false

					if err != nil || result == nil {
						actionResult.Msg = "Error while creating new connection IDa:" + userIDa + " and IDb:" + userIDb
						actionResult.Connected = false
						actionResult.Error = true
					}
				}
			}

		} else {
			actionResult.Msg = "User does not exists"
			actionResult.Connected = false
			actionResult.Error = true

			return actionResult, nil
		}

		return actionResult, nil
	})

	if result == nil {
		return &pb.AddConnectionResult{Msg: "Error"}, err
	} else {
		return result.(*pb.AddConnectionResult), err
	}
}

func (store *ConnectionDBStore) ApproveConnection(userIDa, userIDb string) (*pb.ActionResult, error) {
	actionResult := &pb.ActionResult{Msg: "msg"}
	actionResult.Msg = "Odobravanje konekcije"

	if userIDa == userIDb {
		return &pb.ActionResult{Msg: "UserIDa is same as userIDb"}, nil
	}

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		actionResult := &pb.ActionResult{Msg: "msg"}

		if checkIfUserExist(userIDa, transaction) && checkIfUserExist(userIDb, transaction) {

			//prebacuje status zahtjeva na true -> approved
			_, err := transaction.Run(
				"MATCH (n1{userID:$u1ID})-[r:FRIEND]->(n2{userID:$u2ID}) set r.isApproved = $isApproved RETURN r",
				map[string]interface{}{"u1ID": userIDa, "u2ID": userIDb, "isApproved": true})

			if err != nil {
				actionResult.Msg = "Error while approving connection request with ID:" + userIDb
				return actionResult, err
			}

			//kreira konekciju od user2 do user1
			//TODO:azurirati vrijeme konekcije u1->u2 kad se odobri
			dateNow := time.Now().Local().Unix()
			_, err2 := transaction.Run(
				"MATCH (u1:USER) WHERE u1.userID=$u1ID MATCH (u2:USER) WHERE u2.userID=$u2ID CREATE (u2)-[f:FRIEND{date: $dateNow, isApproved:$isApproved}]->(u1) RETURN u1, u2",
				map[string]interface{}{"u1ID": userIDb, "u2ID": userIDa, "isApproved": true, "dateNow": dateNow})

			if err2 != nil {
				actionResult.Msg = "Error while approving connection request with ID:" + userIDb
				return actionResult, err2
			}

		} else {
			actionResult.Msg = "User does not exist"
			return actionResult, nil
		}

		actionResult.Msg = "Successfully approved connection request IDa:" + userIDa + " and IDb:" + userIDb

		return actionResult, nil
	})

	if result == nil {
		return &pb.ActionResult{Msg: "error"}, err
	} else {
		return result.(*pb.ActionResult), err
	}
}

func (store *ConnectionDBStore) RejectConnection(userIDa, userIDb string) (*pb.ActionResult, error) {
	actionResult := &pb.ActionResult{Msg: "msg"}
	actionResult.Msg = "Odbijanje konekcije"

	if userIDa == userIDb {
		return &pb.ActionResult{Msg: "UserIDa is same as userIDb"}, nil
	}

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		actionResult := &pb.ActionResult{Msg: "msg"}

		if checkIfUserExist(userIDa, transaction) && checkIfUserExist(userIDb, transaction) {
			//TODO:provjeri da li postoji uopste zahtjev/konekcija

			//brise vezu/zahjev
			_, err := transaction.Run(
				"MATCH (u1:USER{userID:$u1ID})<-[rel:FRIEND]-(u2:USER{userID:$u2ID}) DELETE rel",
				map[string]interface{}{"u1ID": userIDa, "u2ID": userIDb})

			if err != nil {
				actionResult.Msg = "Error while rejecting new node with ID:" + userIDb
				return actionResult, err
			}

			//prebrojava broj preostalih veza kod cvorova, ako je 0, obrisacemo cvorove
			result, _ := transaction.Run(
				"MATCH (n:USER{userID:$u1ID})-[rel:FRIEND]-() RETURN COUNT (rel) as broj",
				map[string]interface{}{"u1ID": userIDa})

			//broj veza za userA
			for result.Next() {
				record := result.Record()
				numRelA, _ := record.Get("broj")
				fmt.Println(numRelA)

				if numRelA.(int64) == 0 {
					_, error := transaction.Run(
						"MATCH (u1:USER{userID:$u1ID}) DELETE u1",
						map[string]interface{}{"u1ID": userIDa})

					if error != nil {
						actionResult.Msg = "Error while deleting node with ID:" + userIDa
						return actionResult, err
					}
				}
			}
			resultB, _ := transaction.Run(
				"MATCH (n:USER{userID:$u1ID})-[rel:FRIEND]-() RETURN COUNT (rel) as numRel",
				map[string]interface{}{"u1ID": userIDb})

			//broj veza za userB
			for resultB.Next() {
				record := resultB.Record()
				numRelB, _ := record.Get("numRel")
				fmt.Println(numRelB.(int64))

				if numRelB.(int64) == 0 {
					_, err := transaction.Run(
						"MATCH (u:USER{userID:$u1ID}) DELETE u",
						map[string]interface{}{"u1ID": userIDb})

					if err != nil {
						actionResult.Msg = "Error while deleting node with ID:" + userIDb
						return actionResult, err
					}
				}
			}

		} else {
			actionResult.Msg = "User does not exist"
			return actionResult, nil
		}

		actionResult.Msg = "Successfully rejected connection request IDa:" + userIDa + " to IDb:" + userIDb

		return actionResult, nil
	})

	if result == nil {
		return &pb.ActionResult{Msg: "error"}, err
	} else {
		return result.(*pb.ActionResult), err
	}

}

func (store *ConnectionDBStore) CheckConnection(userIDa, userIDb string) (*pb.ConnectedResult, error) {
	fmt.Println(userIDa)
	fmt.Println(userIDb)

	actionResult := &pb.ConnectedResult{}
	actionResult.Connected = false

	if userIDa == userIDb {
		return &pb.ConnectedResult{Connected: false}, nil
	}

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		actionResult := &pb.ConnectedResult{}

		if checkIfUserExist(userIDa, transaction) && checkIfUserExist(userIDb, transaction) {
			if checkIfFriendExist(userIDa, userIDb, transaction) && checkIfFriendExist(userIDb, userIDa, transaction) {
				actionResult.Connected = true
				actionResult.Request = false
				return actionResult, nil
			}
			if checkIfFriendExist(userIDa, userIDb, transaction) && !checkIfFriendExist(userIDb, userIDa, transaction) {
				actionResult.Connected = false
				actionResult.Request = true
				actionResult.Blocked = false
				return actionResult, nil
			}
			if checkIfFriendExist(userIDb, userIDa, transaction) && !checkIfFriendExist(userIDa, userIDb, transaction) {
				actionResult.Connected = false
				actionResult.Request = true
				actionResult.Blocked = false
				return actionResult, nil
			}

			if checkIfBlockExist(userIDb, userIDa, transaction) || checkIfBlockExist(userIDa, userIDb, transaction) {
				actionResult.Connected = false
				actionResult.Request = false
				actionResult.Blocked = true
				return actionResult, nil
			}

			actionResult.Connected = false
			actionResult.Request = false
			actionResult.Blocked = false

			return actionResult, nil

		} else {
			actionResult.Connected = false
			return actionResult, nil
		}
	})

	if result == nil {
		return &pb.ConnectedResult{Connected: false}, err
	} else {
		return result.(*pb.ConnectedResult), err
	}

}

func (store *ConnectionDBStore) BlockUser(userIDa, userIDb string, isPublic bool, isPublicLogged bool) (*pb.ActionResult, error) {
	actionResult := &pb.ActionResult{Msg: "msg"}
	actionResult.Msg = "Blokiranje korisnika"

	if userIDa == userIDb {
		return &pb.ActionResult{Msg: "UserIDa is same as userIDb"}, nil
	}

	session := (*store.connectionDB).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		actionResult := &pb.ActionResult{Msg: "msg"}

		//ako ne postoji userA, kreira ga
		if !checkIfUserExist(userIDa, transaction) {
			_, err := transaction.Run(
				"CREATE (new_user:USER{userID:$userID, isPublic:$isPublic})",
				map[string]interface{}{"userID": userIDa, "isPublic": isPublicLogged})

			if err != nil {
				actionResult.Msg = "Error while creating new user node with ID:" + userIDa
				return actionResult, err
			}
		}
		//ako ne postoji userB, kreira ga
		if !checkIfUserExist(userIDb, transaction) {
			_, err := transaction.Run(
				"CREATE (new_user:USER{userID:$userID, isPublic:$isPublic})",
				map[string]interface{}{"userID": userIDb, "isPublic": isPublic})

			if err != nil {
				actionResult.Msg = "Error while creating new user node with ID:" + userIDb
				return actionResult, err
			}
		}

		if checkIfUserExist(userIDa, transaction) && checkIfUserExist(userIDb, transaction) {

			//ako je jedan od usera blokirao drugog
			if checkIfBlockExist(userIDa, userIDb, transaction) || checkIfBlockExist(userIDb, userIDa, transaction) {
				fmt.Println("BLOCKED!!!!!")
				actionResult.Msg = "Already blocked" + userIDb
				return actionResult, nil
			}

			if checkIfFriendExist(userIDa, userIDb, transaction) {
				fmt.Println("VEZA 1!!!!!")

				//brise vezu izmedju A i B
				_, err := transaction.Run(
					"MATCH (u1:USER{userID:$u2ID})<-[rel:FRIEND]-(u2:USER{userID:$u1ID}) DELETE rel",
					map[string]interface{}{"u1ID": userIDa, "u2ID": userIDb})

				if err != nil {
					actionResult.Msg = "Error while deleting relationship between ID:" + userIDa + "and ID" + userIDb
					return actionResult, err
				}
			}

			if checkIfFriendExist(userIDb, userIDa, transaction) {
				fmt.Println("VEZA 2!!!!!")

				//brise vezu/zahjev izmedju B i A
				_, err := transaction.Run(
					"MATCH (u1:USER{userID:$u1ID})<-[rel:FRIEND]-(u2:USER{userID:$u2ID}) DELETE rel",
					map[string]interface{}{"u1ID": userIDa, "u2ID": userIDb})

				if err != nil {
					actionResult.Msg = "Error while deleting relationship between ID:" + userIDa + "and ID" + userIDb
					return actionResult, err
				}
			}

			//kreira vezu BLOCK
			result, err := transaction.Run(
				"MATCH (u1:USER) WHERE u1.userID=$uIDa "+
					"MATCH (u2:USER) WHERE u2.userID=$uIDb "+
					"CREATE (u1)-[r1:BLOCK]->(u2) "+
					"RETURN r1", map[string]interface{}{"uIDa": userIDa, "uIDb": userIDb})

			if err != nil && result != nil {
				actionResult.Msg = "Error while blocking request with ID:" + userIDb
				return actionResult, err
			}

		} else {
			actionResult.Msg = "User does not exist"
			return actionResult, nil
		}

		actionResult.Msg = "Successfully blocked IDa:" + userIDa + " and IDb:" + userIDb

		return actionResult, nil
	})

	if result == nil {
		return &pb.ActionResult{Msg: "error"}, err
	} else {
		return result.(*pb.ActionResult), err
	}
}
