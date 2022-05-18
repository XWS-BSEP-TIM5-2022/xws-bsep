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

func (handler *AuthHandler) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return handler.service.Register(ctx, request)
}

func (handler *AuthHandler) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	return handler.service.Login(ctx, request)
}

func (handler *AuthHandler) GetAll(ctx context.Context, request *pb.Empty) (*pb.GetAllResponse, error) {
	return handler.service.GetAll(ctx, request)
}

func (handler *AuthHandler) UpdateUsername(ctx context.Context, request *pb.UpdateUsernameRequest) (*pb.UpdateUsernameResponse, error) {
	return handler.service.UpdateUsername(ctx, request)
}

func (handler *AuthHandler) ChangePassword(ctx context.Context, request *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	return handler.service.ChangePassword(ctx, request)
}

func (handler *AuthHandler) ActivateAccount(ctx context.Context, request *pb.ActivationRequest) (*pb.ActivationResponse, error) {
	return handler.service.ActivateAccount(ctx, request)
}
