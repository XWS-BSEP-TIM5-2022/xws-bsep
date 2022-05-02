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

func (handler *UserHandler) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResponse, error) {

	user := mapInsertUser(request.User)
	success, err := handler.service.Insert(user)
	response := &pb.InsertResponse{
		Success: success,
	}
	return response, err
}

func (handler *UserHandler) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	user := mapInsertUser(request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	return response, err
}

func (handler *UserHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := request.Id
	user, err := handler.service.Get(id)
	if err != nil {
		return nil, err
	}
	userPb := mapUser(user)
	response := &pb.GetResponse{
		User: userPb,
	}
	return response, nil
}
