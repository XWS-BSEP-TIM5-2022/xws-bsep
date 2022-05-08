package application

import (
	"context"
	"fmt"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/infrastructure/persistence"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	store             *persistence.AuthPostgresStore
	jwtService        *JWTService
	userServiceClient user.UserServiceClient
}

func NewAuthService(store *persistence.AuthPostgresStore, jwtService *JWTService, userServiceClient user.UserServiceClient) *AuthService {
	return &AuthService{
		store:             store,
		jwtService:        jwtService,
		userServiceClient: userServiceClient,
	}
}

func (service *AuthService) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userRequest := &user.User{
		Username:     request.Username,
		Name:         request.Name,
		LastName:     request.LastName,
		MobileNumber: request.MobileNumber,
		Gender:       user.User_GenderEnum(request.Gender), // ?
		Birthday:     request.Birthday,
		Email:        request.Email,
		Biography:    request.Biography,
		IsPublic:     request.IsPublic,
	}
	createUserRequest := &user.InsertRequest{
		User: userRequest,
	}
	fmt.Println(createUserRequest)

	createUserResponse, err := service.userServiceClient.Insert(context.TODO(), createUserRequest)
	if err != nil {
		return nil, err
	}
	authCredentials, err := domain.NewAuthCredentials(
		createUserResponse.Id,
		request.Username,
		request.Password,
	)
	if err != nil {
		return nil, err
	}
	authCredentials, err = service.store.Create(authCredentials)
	if err != nil {
		return nil, err
	} else {
		token, err := service.jwtService.GenerateToken(authCredentials)
		if err != nil {
			return nil, err
		}
		return &pb.RegisterResponse{
			Token: token,
		}, nil
	}
}

func (service *AuthService) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	authCredentials, err := service.store.FindByUsername(request.Username)
	if err != nil {
		return nil, err
	}
	fmt.Println("No error finding user")

	ok := authCredentials.CheckPassword(request.Password)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid username or password")
	}
	fmt.Println("No error validating password")
	token, err := service.jwtService.GenerateToken(authCredentials)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not generate JWT token")
	}
	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (service *AuthService) GetAll(ctx context.Context, request *pb.Empty) (*pb.GetAllResponse, error) {
	auths, err := service.store.FindAll()
	if err != nil || *auths == nil {
		return nil, err
	}
	response := &pb.GetAllResponse{
		Auth: []*pb.Auth{},
	}
	for _, auth := range *auths {
		current := &pb.Auth{
			Id:       auth.Id,
			Username: auth.Username,
			Password: auth.Password,
			Role:     auth.Role,
		}
		response.Auth = append(response.Auth, current)
	}
	return response, nil
}
