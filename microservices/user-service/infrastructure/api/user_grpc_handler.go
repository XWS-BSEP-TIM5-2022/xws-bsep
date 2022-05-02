package api

import (
	"context"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/application"
	domain "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
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

	emptUser := domain.User{}
	emptUser.Id = request.User.Id
	emptUser.Name = request.User.Name
	// user := mapUser(request.User)
	success, err := handler.service.Insert(&emptUser)
	response := &pb.InsertResponse{
		Success: success,
	}
	return response, err
}

//func (handler *UserHandler) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {
//	user := mapUser(request.User)
//	success, err := handler.service.UpdateAllInfo(user)
//	response := &pb.UpdateResponse{
//		Success: success,
//	}
//	return response, err
//}
