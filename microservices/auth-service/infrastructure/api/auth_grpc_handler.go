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

func (handler *AuthHandler) PasswordlessLogin(ctx context.Context, request *pb.PasswordlessLoginRequest) (*pb.PasswordlessLoginResponse, error) {
	return handler.service.PasswordlessLogin(ctx, request)
}

func (handler *AuthHandler) ConfirmEmailLogin(ctx context.Context, request *pb.ConfirmEmailLoginRequest) (*pb.ConfirmEmailLoginResponse, error) {
	return handler.service.ConfirmEmailLogin(ctx, request)
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

func (handler *AuthHandler) SendRecoveryCode(ctx context.Context, request *pb.SendRecoveryCodeRequest) (*pb.SendRecoveryCodeResponse, error) {
	return handler.service.SendRecoveryCode(ctx, request)
}

func (handler *AuthHandler) VerifyRecoveryCode(ctx context.Context, request *pb.VerifyRecoveryCodeRequest) (*pb.Response, error) {
	return handler.service.VerifyRecoveryCode(ctx, request)
}

func (handler *AuthHandler) ResetForgottenPassword(ctx context.Context, request *pb.ResetForgottenPasswordRequest) (*pb.Response, error) {
	return handler.service.ResetForgottenPassword(ctx, request)
}

func (handler *AuthHandler) GetAllPermissionsByRole(ctx context.Context, request *pb.Empty) (*pb.Response, error) {
	return handler.service.GetAllPermissionsByRole(ctx, request)
}

func (handler *AuthHandler) AdminsEndpoint(ctx context.Context, request *pb.Empty) (*pb.Response, error) {
	return handler.service.AdminsEndpoint(ctx, request)
}

func (handler *AuthHandler) CreateNewAPIToken(ctx context.Context, request *pb.APITokenRequest) (*pb.NewAPITokenResponse, error) {
	return handler.service.CreateNewAPIToken(ctx, request)
}
