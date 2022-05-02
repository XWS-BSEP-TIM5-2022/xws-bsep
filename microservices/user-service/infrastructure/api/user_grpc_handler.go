package api

import (
	"context"

	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/application"
)

type UserHandler struct {
	service *application.UserService
	pb.UnimplementedUserServiceServer
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (handler *UserHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	users, err := handler.service.GetAll()
	if err != nil || *users == nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}
	for _, user := range *users {
		current := mapUser(&user)
		response.Users = append(response.Users, current)
	}
	return response, nil
}

// func (handler *UserHandler) Insert(ctx context.Context, request *pb.CreateUserRequest) *pb.CreateUserResponse {
// 	// user := mapUser(&request.User)
// 	// return handler.service.Create(&user)
// 	response := &pb.CreateUserResponse{
// 		User: *pb.User{},
// 	}

// }

func (handler *UserHandler) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResponse, error) {
	user := mapUser(request.User)
	success, err := handler.service.Insert(user)
	response := &pb.InsertResponse{
		Success: success,
	}
	return response, err
}
