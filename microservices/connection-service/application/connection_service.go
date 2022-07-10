package application

import (
	"context"

	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
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

func (service *ConnectionService) GetConnections(ctx context.Context, id string) ([]*domain.UserConn, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConnections service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	var friendsRetVal []*domain.UserConn

	friends, err := service.store.GetConnections(ctx, id)
	if err != nil {
		return nil, nil
	}
	for _, s := range friends {
		friendsRetVal = append(friendsRetVal, &domain.UserConn{UserID: s.UserID, IsPublic: s.IsPublic})
	}
	return friendsRetVal, nil
}

func (service *ConnectionService) GetRequests(ctx context.Context, id string) ([]*domain.UserConn, error) {
	span := tracer.StartSpanFromContext(ctx, "GetRequests service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	var friendsRetVal []*domain.UserConn

	friends, err := service.store.GetRequests(ctx, id)
	if err != nil {
		return nil, nil
	}
	for _, s := range friends {
		friendsRetVal = append(friendsRetVal, &domain.UserConn{UserID: s.UserID, IsPublic: s.IsPublic})
	}
	return friendsRetVal, nil
}

func (service *ConnectionService) Register(ctx context.Context, userID string, isPublic bool) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "GetRequests service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.Register(ctx, userID, isPublic)
}

func (service *ConnectionService) AddConnection(ctx context.Context, userIDa string, userIDb string, isPublic bool, isPublicLogged bool) (*pb.AddConnectionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "AddConnection service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.AddConnection(ctx, userIDa, userIDb, isPublic, isPublicLogged)
}

func (service *ConnectionService) BlockUser(ctx context.Context, userIDa, userIDb string, isPublic bool, isPublicLogged bool) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "BlockUser service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.BlockUser(ctx, userIDa, userIDb, isPublic, isPublicLogged)
}

func (service *ConnectionService) ApproveConnection(ctx context.Context, userIDa, userIDb string) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "ApproveConnection service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.ApproveConnection(ctx, userIDa, userIDb)
}

func (service *ConnectionService) RejectConnection(ctx context.Context, userIDa, userIDb string) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "RejectConnection service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.RejectConnection(ctx, userIDa, userIDb)
}

func (service *ConnectionService) CheckConnection(ctx context.Context, userIDa, userIDb string) (*pb.ConnectedResult, error) {
	span := tracer.StartSpanFromContext(ctx, "CheckConnection service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.CheckConnection(ctx, userIDa, userIDb)
}

func (service *ConnectionService) GetRecommendation(ctx context.Context, userID string) ([]*domain.UserConn, error) {
	span := tracer.StartSpanFromContext(ctx, "GetRecommendation service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.GetRecommendation(ctx, userID)
}

func (service *ConnectionService) ChangePrivacy(ctx context.Context, userIDa string, isPrivate bool) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "ChangePrivacy service")
	defer span.Finish()
	ctx = tracer.ContextWithSpan(context.Background(), span)

	return service.store.ChangePrivacy(ctx, userIDa, isPrivate)
}
