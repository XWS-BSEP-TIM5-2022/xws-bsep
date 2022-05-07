package domain

import pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"

type ConnectionStore interface {
	GetFriends(id string) ([]UserConn, error)
	AddConnection(userIDa, userIDb string) (*pb.ActionResult, error)
	Register(userID string, isPublic bool) (*pb.ActionResult, error)
	ApproveConnection(userIDa, userIDb string) (*pb.ActionResult, error)
}
