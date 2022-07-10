package application

import (
	"context"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/infrastructure/api"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
)

type AuthHandler struct {
	service *api.AuthService
	pb.UnimplementedAuthServiceServer
}

func NewAuthHandler(service *api.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (handler *AuthHandler) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "Login")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.Login(ctx, request)
}

func (handler *AuthHandler) PasswordlessLogin(ctx context.Context, request *pb.PasswordlessLoginRequest) (*pb.PasswordlessLoginResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "PasswordlessLogin")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.PasswordlessLogin(ctx, request)
}

func (handler *AuthHandler) ConfirmEmailLogin(ctx context.Context, request *pb.ConfirmEmailLoginRequest) (*pb.ConfirmEmailLoginResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "ConfirmEmailLogin")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.ConfirmEmailLogin(ctx, request)
}

func (handler *AuthHandler) GetAll(ctx context.Context, request *pb.Empty) (*pb.GetAllResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAll")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.GetAll(ctx, request)
}

func (handler *AuthHandler) UpdateUsername(ctx context.Context, request *pb.UpdateUsernameRequest) (*pb.UpdateUsernameResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "UpdateUsername")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.UpdateUsername(ctx, request)
}

func (handler *AuthHandler) ChangePassword(ctx context.Context, request *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "ChangePassword")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.ChangePassword(ctx, request)
}

func (handler *AuthHandler) ActivateAccount(ctx context.Context, request *pb.ActivationRequest) (*pb.ActivationResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "ActivateAccount")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.ActivateAccount(ctx, request)
}

func (handler *AuthHandler) SendRecoveryCode(ctx context.Context, request *pb.SendRecoveryCodeRequest) (*pb.SendRecoveryCodeResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "SendRecoveryCode")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.SendRecoveryCode(ctx, request)
}

func (handler *AuthHandler) VerifyRecoveryCode(ctx context.Context, request *pb.VerifyRecoveryCodeRequest) (*pb.Response, error) {
	span := tracer.StartSpanFromContext(ctx, "VerifyRecoveryCode")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.VerifyRecoveryCode(ctx, request)
}

func (handler *AuthHandler) ResetForgottenPassword(ctx context.Context, request *pb.ResetForgottenPasswordRequest) (*pb.Response, error) {
	span := tracer.StartSpanFromContext(ctx, "ResetForgottenPassword")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.ResetForgottenPassword(ctx, request)
}

func (handler *AuthHandler) GetAllPermissionsByRole(ctx context.Context, request *pb.Empty) (*pb.Response, error) {
	span := tracer.StartSpanFromContext(ctx, "GetAllPermissionsByRole")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.GetAllPermissionsByRole(ctx, request)
}

func (handler *AuthHandler) AdminsEndpoint(ctx context.Context, request *pb.Empty) (*pb.Response, error) {
	span := tracer.StartSpanFromContext(ctx, "AdminsEndpoint")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.AdminsEndpoint(ctx, request)
}

func (handler *AuthHandler) CreateNewAPIToken(ctx context.Context, request *pb.APITokenRequest) (*pb.NewAPITokenResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateNewAPIToken")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.CreateNewAPIToken(ctx, request)
}

func (handler *AuthHandler) GetUsernameByApiToken(ctx context.Context, request *pb.GetUsernameRequest) (*pb.GetUsernameResponse, error) {
	span := tracer.StartSpanFromContext(ctx, "GetUsernameByApiToken")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)
	return handler.service.GetUsernameByApiToken(ctx, request)
}

func (handler *AuthHandler) GetAllEvents(ctx context.Context, request *pb.GetAllEventsRequest) (*pb.GetAllEventsResponse, error) {

	events, err := handler.service.GetAllEvents()

	if err != nil {
		return nil, err
	}

	var finalEvents []*pb.Event

	for _, event := range events {
		finalEvents = append(finalEvents, mapEvent(event))
	}

	response := &pb.GetAllEventsResponse{
		Events: finalEvents,
	}

	return response, nil

}

func mapEvent(event *domain.Event) *pb.Event {
	eventPb := &pb.Event{
		Id:     event.Id.Hex(),
		UserId: event.UserId,
		Text:   event.Text,
		Date:   event.Date.String(),
	}

	return eventPb
}
