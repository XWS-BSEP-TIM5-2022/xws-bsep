package domain

import (
	"context"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
)

type ConnectionStore interface {
	GetConnections(ctx context.Context, id string) ([]UserConn, error)
	AddConnection(ctx context.Context, userIDa string, userIDb string, isPublic bool, isPublicLogged bool) (*pb.AddConnectionResult, error)
	BlockUser(ctx context.Context, userIDa, userIDb string, isPublic bool, isPublicLogged bool) (*pb.ActionResult, error)
	Register(ctx context.Context, userID string, isPublic bool) (*pb.ActionResult, error)
	ApproveConnection(ctx context.Context, userIDa, userIDb string) (*pb.ActionResult, error)
	RejectConnection(ctx context.Context, userIDa, userIDb string) (*pb.ActionResult, error)
	CheckConnection(ctx context.Context, userIDa, userIDb string) (*pb.ConnectedResult, error)
	GetRequests(ctx context.Context, id string) ([]UserConn, error)
	GetRecommendation(ctx context.Context, userID string) ([]*UserConn, error)
	ChangePrivacy(ctx context.Context, userIDa string, isPrivate bool) (*pb.ActionResult, error)
}
