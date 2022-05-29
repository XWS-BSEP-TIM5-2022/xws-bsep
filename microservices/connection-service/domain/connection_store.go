package domain

import pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"

type ConnectionStore interface {
	GetConnections(id string) ([]UserConn, error)
	AddConnection(userIDa string, userIDb string, isPublic bool) (*pb.AddConnectionResult, error)
	Register(userID string, isPublic bool) (*pb.ActionResult, error)
	ApproveConnection(userIDa, userIDb string) (*pb.ActionResult, error)
	RejectConnection(userIDa, userIDb string) (*pb.ActionResult, error)
	CheckConnection(userIDa, userIDb string) (*pb.ConnectedResult, error)
	GetRequests(id string) ([]UserConn, error)
}
