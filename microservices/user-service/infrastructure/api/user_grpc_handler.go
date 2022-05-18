package api

import (
	"context"
	"fmt"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
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
	user, err := handler.service.Insert(user)

	if err != nil {
		return nil, err
	} else {
		return &pb.InsertResponse{
			Id: user.Id.Hex(),
		}, nil
	}
}

func (handler *UserHandler) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

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

func (handler *UserHandler) GetLoggedInUserInfo(ctx context.Context, request *pb.GetAllRequest) (*pb.User, error) {
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	user, err := handler.service.GetById(userId)
	if err != nil {
		return nil, err
	}
	pbUser := mapUser(user)
	return pbUser, nil
}

func (handler *UserHandler) UpdateBasicInfo(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapBasicInfo(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	return response, err
}

func (handler *UserHandler) UpdateExperienceAndEducation(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapExperienceAndEducation(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	return response, err
}

func (handler *UserHandler) UpdateSkillsAndInterests(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapSkillsAndInterests(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	return response, err
}
 
func (handler *UserHandler) GetEmail(ctx context.Context, request *pb.GetRequest) (*pb.GetEmailResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	user, err := handler.service.Get(objectId)
	if err != nil {
		return nil, err
	}

	response := &pb.GetEmailResponse{
		Email: user.Email,
	}
	return response, nil
}
func (handler *UserHandler) UpdateIsActiveById(ctx context.Context, request *pb.ActivateAccountRequest) (*pb.ActivateAccountResponse, error) {
	err := handler.service.UpdateIsActiveById(request.Id)
	if err != nil {
		return &pb.ActivateAccountResponse{
			Success: err.Error(),
		}, err
	}
	return &pb.ActivateAccountResponse{
		Success: "Success",
	}, nil
}

func (handler *UserHandler) GetIsActive(ctx context.Context, request *pb.GetRequest) (*pb.IsActiveResponse, error) {
	fmt.Println(request.Id)
	user, err := handler.service.GetById(request.Id)
	if err != nil {
		fmt.Println("* error :", err)
		return nil, err
	}
	return &pb.IsActiveResponse{
		IsActive: user.IsActive,
	}, nil 
}
