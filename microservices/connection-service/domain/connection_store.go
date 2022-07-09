package domain

import pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"

type ConnectionStore interface {
	GetConnections(id string) ([]UserConn, error)
	AddConnection(userIDa string, userIDb string, isPublic bool, isPublicLogged bool) (*pb.AddConnectionResult, error)
	BlockUser(userIDa, userIDb string, isPublic bool, isPublicLogged bool) (*pb.ActionResult, error)
	Register(userID string, isPublic bool) (*pb.ActionResult, error)
	ApproveConnection(userIDa, userIDb string) (*pb.ActionResult, error)
	RejectConnection(userIDa, userIDb string) (*pb.ActionResult, error)
	CheckConnection(userIDa, userIDb string) (*pb.ConnectedResult, error)
	GetRequests(id string) ([]UserConn, error)
	GetRecommendation(userID string) ([]*UserConn, error)
	ChangePrivacy(userIDa string, isPrivate bool) (*pb.ActionResult, error)
}
