package api

import (
	"context"
	"fmt"
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

func (handler *ConnectionHandler) GetFriends(ctx context.Context, request *pb.GetRequest) (*pb.Users, error) {

	fmt.Println("[ConnectionHandler]:GetFriends")

	id := request.UserID
	friends, err := handler.service.GetFriends(id)
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
	userIDa := request.AddConnectionDTO.UserIDa
	userIDb := request.AddConnectionDTO.UserIDb
	return handler.service.AddConnection(userIDa, userIDb)
}
