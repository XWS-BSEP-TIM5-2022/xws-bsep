package application

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/connection_service/domain"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
)

type ConnectionService struct {
	store domain.ConnectionStore
}

func NewConnectionService(store domain.ConnectionStore) *ConnectionService {
	return &ConnectionService{
		store: store,
	}
}

func (service *ConnectionService) GetFriends(id string) ([]*domain.UserConn, error) {

	var friendsRetVal []*domain.UserConn

	friends, err := service.store.GetFriends(id)
	if err != nil {
		return nil, nil
	}
	for _, s := range friends {
		friendsRetVal = append(friendsRetVal, &domain.UserConn{UserID: s.UserID, IsPublic: s.IsPublic})
	}
	return friendsRetVal, nil
}

func (service *ConnectionService) AddFriend(userIDa, userIDb string) (*pb.ActionResult, error) {
	return service.store.AddFriend(userIDa, userIDb), nil
}
