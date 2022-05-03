package domain

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
)

type ConnectionStore interface {
	GetFriends(id string) ([]UserConn, error)
	AddFriend(userIDa, userIDb string) (*pb.ActionResult, error)
}
