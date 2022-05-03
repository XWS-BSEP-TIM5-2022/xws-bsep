package api

import (
	"auth-service/application"
	"context"
	"fmt"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
)

type AuthHandler struct {
	service *application.AuthService
	pb.UnimplementedAuthServiceServer
}

func NewAuthHandler(service *application.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// func (handler *AuthHandler) GetAll(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
// 	auths, err := handler.service.GetAll()
// 	if err != nil || *auths == nil {
// 		return nil, err
// 	}
// 	response := &pb.GetResponse{
// 		Auth: *pb.Authentication{},
// 	}
// 	for _, user := range *users {
// 		current := mapUser(&user)
// 		response.Users = append(response.Users, current)
// 	}
// 	return response, nil
// }

func (handler *AuthHandler) Insert(ctx context.Context, request *pb.AddRequest) (*pb.AddResponse, error) {
	authentication := domain.Authentication{}
	authentication.Id = request.Auth.Id
	authentication.Name = request.Auth.Username
	authentication.Password = request.Auth.Password
	authentication.Role = request.Auth.Role

	success, err := handler.service.Create(&authentication)
	if err != nil {
		success := "Greska prilikom upisa u bazu!"
		response := &pb.AddResponse{
			Success: success,
		}
		return response, err
	}
	fmt.Println(success)
	response := &pb.AddResponse{
		Success: success,
	}
	return response, nil
}
