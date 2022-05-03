package api

import (
	"context"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/application"
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

func (handler *AuthHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	auths, err := handler.service.GetAll()
	if err != nil || *auths == nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Authentications: []*pb.Auth{},
	}
	for _, auth := range *auths {
		current := pb.Auth{
			Id:       auth.Id,
			Name:     auth.Name,
			Password: auth.Password,
		}
		response.Authentications = append(response.Authentications, &current)
	}
	return response, nil
}

// func (handler *AuthHandler) Create(ctx context.Context, request *pb.AddRequest) (*pb.AddResponse, error) {

// 	success := "TODO: zavrsiti mapiranje"
// 	// auth := mapCreateAuth(request.Auth)
// 	// fmt.Println(auth)
// 	// success, err := handler.service.Create(&current)
// 	// if err != nil {
// 	// 	success := "Greska prilikom upisa u bazu!"
// 	// 	response := &pb.AddResponse{
// 	// 		Success: success,
// 	// 	}
// 	// 	return response, err
// 	// }
// 	fmt.Println(success)
// 	response := &pb.AddResponse{
// 		Success: success,
// 	}
// 	// return response, err
// 	return response, nil
// }
