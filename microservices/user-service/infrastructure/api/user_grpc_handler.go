package api

import (
	"context"
	"errors"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	return response, nil
}

func (handler *UserHandler) GetAllPublic(ctx context.Context, request *pb.GetAllPublicRequest) (*pb.GetAllPublicResponse, error) {
	users, err := handler.service.GetAllPublic()
	if err != nil {
		return nil, err
	}
	response := &pb.GetAllPublicResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	return response, nil
}

func (handler *UserHandler) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {

	criteria := request.Criteria
	users, err := handler.service.Search(criteria)

	if err != nil {
		return nil, err
	}

	response := &pb.SearchResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}

	return response, nil
}
func (handler *UserHandler) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResponse, error) {

	user := mapInsertUser(request.User)
	success, err := handler.service.Insert(user)

	if err != nil {
		return nil, err
	}

	response := &pb.InsertResponse{
		Success: success,
	}
	return response, err
}

func (handler *UserHandler) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id, _ := primitive.ObjectIDFromHex(request.User.Id)
	oldUser, err := handler.service.Get(id)

	if err != nil {
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapUpdateUser(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	return response, err
}

func (handler *UserHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {

	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}
	userPb := mapUser(user)
	response := &pb.GetResponse{
		User: userPb,
	}
	return response, nil
}

func (handler *UserHandler) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {

	user, err := handler.service.GetByUsername(request.GetData().Username)
	if err != nil {
		return &pb.LoginResponse{
			Success: "there is no user with that username",
		}, errors.New("there is no user with that username")
	}

	if request.GetData().Password != user.Password {
		return &pb.LoginResponse{
			Success: "passwords do not match",
		}, errors.New("passwords do not match")
	}

	return &pb.LoginResponse{
		Success: "success",
	}, nil
}
