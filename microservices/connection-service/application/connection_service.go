package application

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/domain"
)

type ConnectionService struct {
	store domain.ConnectionStore
}

func NewConnectionService(store domain.ConnectionStore) *ConnectionService {
	return &ConnectionService{
		store: store,
	}
}

func (service *ConnectionService) GetConnections(id string) ([]*domain.UserConn, error) {

	var friendsRetVal []*domain.UserConn

	friends, err := service.store.GetConnections(id)
	if err != nil {
		return nil, nil
	}
	for _, s := range friends {
		friendsRetVal = append(friendsRetVal, &domain.UserConn{UserID: s.UserID, IsPublic: s.IsPublic})
	}
	return friendsRetVal, nil
}

func (service *ConnectionService) Register(userID string, isPublic bool) (*pb.ActionResult, error) {
	return service.store.Register(userID, isPublic)
}

func (service *ConnectionService) AddConnection(userIDa string, userIDb string, isPublic bool) (*pb.AddConnectionResult, error) {
	return service.store.AddConnection(userIDa, userIDb, isPublic)
}

func (service *ConnectionService) ApproveConnection(userIDa, userIDb string) (*pb.ActionResult, error) {
	return service.store.ApproveConnection(userIDa, userIDb)
}

func (service *ConnectionService) RejectConnection(userIDa, userIDb string) (*pb.ActionResult, error) {
	return service.store.RejectConnection(userIDa, userIDb)
}

func (service *ConnectionService) CheckConnection(userIDa, userIDb string) (*pb.ConnectedResult, error) {
	return service.store.CheckConnection(userIDa, userIDb)
}
