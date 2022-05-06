package application

import (
	"context"
	"fmt"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/infrastructure/persistence/persistence"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
)

type AuthService struct {
	store             *persistence.AuthPostgresStore
	jwtService        *JWTManager
	userServiceClient user.UserServiceClient
}

func NewAuthService(store *persistence.AuthPostgresStore, jwtService *JWTManager, userServiceClient user.UserServiceClient) *AuthService {
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
		// Gender:       request.Gender,
		Birthday:  request.Birthday,
		Email:     request.Email,
		Biography: request.Biography,
		Password:  request.Password,
		IsPublic:  request.IsPublic,
	}
	createUserRequest := &user.InsertRequest{
		User: userRequest,
	}

	fmt.Println(createUserRequest)

	// Register user
	createUserResponse, err := service.userServiceClient.Insert(context.TODO(), createUserRequest)
	if err != nil {
		return nil, err
	}

	// kreiraju se auth kredencijali preko konstruktora da bi mogla odmah da se hesira lozinka
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
