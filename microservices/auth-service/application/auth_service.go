package application

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		// Username:     request.Username,
		Name:         request.Name,
		LastName:     request.LastName,
		MobileNumber: request.MobileNumber,
		Gender:       user.User_GenderEnum(request.Gender), // ?
		Birthday:     request.Birthday,
		Email:        request.Email,
		Biography:    request.Biography,
		IsPublic:     request.IsPublic,
	}

	for _, education := range request.Education {

		ed_id := primitive.NewObjectID().Hex()

		userRequest.Education = append(userRequest.Education, &user.Education{
			Id:        ed_id,
			Name:      education.Name,
			Level:     user.Education_EducationEnum(education.Level),
			Place:     education.Place,
			StartDate: education.StartDate,
			EndDate:   education.EndDate,
		})
	}

	for _, experience := range request.Experience {

		ex_id := primitive.NewObjectID().Hex()

		userRequest.Experience = append(userRequest.Experience, &user.Experience{
			Id:        ex_id,
			Name:      experience.Name,
			Headline:  experience.Headline,
			Place:     experience.Place,
			StartDate: experience.StartDate,
			EndDate:   experience.EndDate,
		})
	}

	for _, skill := range request.Skills {

		s_id := primitive.NewObjectID().Hex()

		userRequest.Skills = append(userRequest.Skills, &user.Skill{
			Id:   s_id,
			Name: skill.Name,
		})
	}

	for _, interest := range request.Interests {

		in_id := primitive.NewObjectID().Hex()

		userRequest.Interests = append(userRequest.Interests, &user.Interest{
			Id:          in_id,
			Name:        interest.Name,
			Description: interest.Description,
		})
	}

	createUserRequest := &user.InsertRequest{
		User: userRequest,
	}

	auths, err := service.store.FindAll()
	for _, auth := range *auths {
		if auth.Username == request.Username {
			log.Println("Username is not unique")
			return nil, errors.New("Username is not unique")
		}
	}

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

func (service *AuthService) UpdateUsername(ctx context.Context, request *pb.UpdateUsernameRequest) (*pb.UpdateUsernameResponse, error) {
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	if userId == "" {
		return &pb.UpdateUsernameResponse{
			StatusCode: "500",
			Message:    "User id not found",
		}, nil
	} else {
		auths, err := service.store.FindAll()
		for _, auth := range *auths {
			if auth.Username == request.Username {
				log.Println("Username is not unique")
				return &pb.UpdateUsernameResponse{
					StatusCode: "500",
					Message:    "Username is not unique",
				}, errors.New("Username is not unique")
			}
		}
		response, err := service.store.UpdateUsername(userId, request.Username)
		if err != nil {
			return &pb.UpdateUsernameResponse{
				StatusCode: "500",
				Message:    "Auth service credentials not found from JWT token",
			}, err
		}
		log.Print(response)
		return &pb.UpdateUsernameResponse{
			StatusCode: "200",
			Message:    "Username updated",
		}, nil
	}

}
