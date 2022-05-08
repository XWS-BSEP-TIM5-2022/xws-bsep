package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/application"
)

type ConnectionHandler struct {
	pb.UnimplementedConnectionServiceServer
	service *application.ConnectionService
}

func NewConnectionHandler(service *application.ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{
		service: service,
	}
}

func (handler *ConnectionHandler) GetConnections(ctx context.Context, request *pb.GetRequest) (*pb.Users, error) {

	fmt.Println("[ConnectionHandler]:GetConnections")

	id := request.UserID
	//prosledili smo registrovanog korisnika
	//id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	friends, err := handler.service.GetConnections(id)
	if err != nil {
		return nil, err
	}
	response := &pb.Users{}
	for _, user := range friends {
		fmt.Println("User", id, "is friend with", user.UserID)
		response.Users = append(response.Users, mapUserConn(user))
	}
	return response, nil
}

func (handler *ConnectionHandler) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.ActionResult, error) {
	fmt.Println("[ConnectionHandler]:Register")
	userID := request.User.UserID
	isPublic := request.User.IsPublic
	return handler.service.Register(userID, isPublic)
}

func (handler *ConnectionHandler) AddConnection(ctx context.Context, request *pb.AddConnectionRequest) (*pb.ActionResult, error) {
	fmt.Println("[ConnectionHandler]:AddConnection")

	//prosledili smo registrovanog korisnika
	userIDa := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	userIDb := request.AddConnectionDTO.UserIDb
	return handler.service.AddConnection(userIDa, userIDb)
}

func (handler *ConnectionHandler) ApproveConnection(ctx context.Context, request *pb.ApproveConnectionRequest) (*pb.ActionResult, error) {
	fmt.Println("[ConnectionHandler]:ApproveConnection")
	//prosledili smo registrovanog korisnika
	userIDa := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	userIDb := request.ApproveConnectionDTO.UserIDb
	return handler.service.ApproveConnection(userIDa, userIDb)
}

func (handler *ConnectionHandler) RejectConnection(ctx context.Context, request *pb.RejectConnectionRequest) (*pb.ActionResult, error) {
	fmt.Println("[ConnectionHandler]:RejectConnection")
	//prosledili smo registrovanog korisnika
	userIDa := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	userIDb := request.RejectConnectionDTO.UserIDb
	return handler.service.RejectConnection(userIDa, userIDb)
}
