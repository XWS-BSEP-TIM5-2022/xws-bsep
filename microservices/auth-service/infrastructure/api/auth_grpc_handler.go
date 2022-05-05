package api

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/application"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			Date:     timestamppb.New(auth.Date),
			Role:     auth.Role,
		}
		response.Authentications = append(response.Authentications, &current)
	}
	return response, nil
}

func (handler *AuthHandler) Create(ctx context.Context, request *pb.AddRequest) (*pb.AddResponse, error) {
	auth := mapCreateAuth(request.Auth)
	tokenString := GenerateToken(&ctx, request.Auth.Id, request.Auth.Name, request.Auth.Role)
	if tokenString == "" {
		success := "Greska prilikom generisanja tokena!"
		response := &pb.AddResponse{
			Success: success,
		}
		return response, nil
	}
	success, err := handler.service.Create(auth)
	if err != nil {
		mess := "Greska prilikom upisa auth kredencijala u bazu!"
		response := &pb.AddResponse{
			Success: mess,
		}
		return response, err
	}
	fmt.Println(success)
	response := &pb.AddResponse{
		Success: tokenString,
	}
	return response, nil
}

func (handler *AuthHandler) ValidateToken(ctx context.Context, request *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	err := ValidateToken(request.Token)
	if err != nil {
		response := &pb.ValidateResponse{
			Id:       "",
			Username: "",
			Role:     "",
		}
		return response, err
	}
	response := &pb.ValidateResponse{
		Id:       "neko",
		Username: "neko",
		Role:     "neko",
	}
	return response, nil
}

//  ---------------------------------- JWT ------------------------------------
var jwtKey = []byte("supersecretkey")

type JWTClaim struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(id string, username string, role string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Minute)
	claims := &JWTClaim{
		Id:       id,
		Role:     role,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	fmt.Println("----------------- CLAIMS -------------------")
	fmt.Println(claims.Id)
	fmt.Println(claims.Username)
	fmt.Println(claims.Role)
	fmt.Println(claims.ExpiresAt)
	fmt.Println("----------------- ------ -------------------")
	fmt.Println(" ******* validate token ******* ")
	fmt.Println(err)
	fmt.Println("************************* ")
	return
}

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GenerateToken(context *context.Context, id, username, role string) string {
	// var request TokenRequest
	// var user models.User
	// if err := context.ShouldBindJSON(&request); err != nil {
	// 	context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 	context.Abort()
	// 	return
	// }
	// check if email exists and password is correct
	// record := database.Instance.Where("email = ?", request.Email).First(&user)
	// if record.Error != nil {
	// 	context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
	// 	context.Abort()
	// 	return
	// }
	// credentialError := user.CheckPassword(request.Password)
	// if credentialError != nil {
	// 	context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	// 	context.Abort()
	// 	return
	// }

	tokenString, err := GenerateJWT(id, username, role)
	if err != nil {
		// TODO: exception
		fmt.Println("GRESKA PRI GENERISANJU TOKENA")
		return ""
	}
	fmt.Println(tokenString)
	return tokenString
	// if err != nil {
	// 	context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	context.Abort()
	// 	return
	// }
	// context.JSON(http.StatusOK, gin.H{"token": tokenString})
}
