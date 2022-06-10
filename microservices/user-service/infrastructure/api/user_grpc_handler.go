package api

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/application"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/peer"
)

type UserHandler struct {
	service      *application.UserService
	CustomLogger *CustomLogger
	pb.UnimplementedUserServiceServer
}

func NewUserHandler(service *application.UserService) *UserHandler {
	// InfoLogger := setLogrusLogger(config.NewConfig().InfoLogsFile)
	// ErrorLogger := setLogrusLogger(config.NewConfig().ErrorLogsFile)
	// WarningLogger := setLogrusLogger(config.NewConfig().WarningLogsFile)
	// SuccessLogger := setLogrusLogger(config.NewConfig().SuccessLogsFile)
	// DebugLogger := setLogrusLogger(config.NewConfig().DebugLogsFile)

	// writer, _ := rotatelogs.New(
	// 	"logs/cao"+".%Y%m%d%H%M"+".log",
	// 	rotatelogs.WithLinkName("logs/cao.log"),
	// 	rotatelogs.WithMaxAge(time.Duration(60)*time.Second),
	// 	rotatelogs.withMaxSize(),
	// 	rotatelogs.WithRotationTime(time.Duration(60)*time.Second),
	// )

	// log.SetOutput(&lumberjack.Logger{
	// 	Filename:   "logsTest/foo.log",
	// 	MaxSize:    1, // megabytes
	// 	MaxBackups: 10,
	// 	MaxAge:     7,    //days
	// 	Compress:   true, // disabled by default
	// })

	CustomLogger := NewCustomLogger()
	return &UserHandler{
		service:      service,
		CustomLogger: CustomLogger,
	}
}

func (handler *UserHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	users, err := handler.service.GetAll()
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Get all")
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(users)) + " public users")
	return response, nil
}

func (handler *UserHandler) GetAllPublic(ctx context.Context, request *pb.GetAllPublicRequest) (*pb.GetAllPublicResponse, error) {
	handler.CustomLogger.InfoLogger.Info("Getting all public accounts") //.InfoLogger.WithFields(logrus.Fields{}).Info(...)
	// SD: host and port from request
	p, _ := peer.FromContext(ctx)
	fmt.Println("** ** IP: " + p.Addr.String())

	users, err := handler.service.GetAllPublic()
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Found " + strconv.Itoa(len(users)) + " public users")
		return nil, err
	}
	response := &pb.GetAllPublicResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(users)) + " public users")
	return response, nil
}

func (handler *UserHandler) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	criteria := request.Criteria
	users, err := handler.service.Search(criteria)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Search error")
		return nil, err
	}

	response := &pb.SearchResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	handler.CustomLogger.SuccessLogger.Info("Number of users found after search: " + strconv.Itoa(len(users)))
	return response, nil
}
func (handler *UserHandler) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResponse, error) {
	user := mapInsertUser(request.User)
	user, err := handler.service.Insert(user)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User is not inserted")
		return nil, err
	} else {
		handler.CustomLogger.SuccessLogger.Info("User inserted successfully")
		return &pb.InsertResponse{
			Id: user.Id.Hex(),
		}, nil
	}
}

func (handler *UserHandler) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	//id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(request.User.Id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created")
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapUpdateUser(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("User with ID: " + user.Id.Hex() + "updated successfully")
	return response, err
}

func (handler *UserHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	user, err := handler.service.Get(objectId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return nil, err
	}
	userPb := mapUser(user)
	response := &pb.GetResponse{
		User: userPb,
	}
	handler.CustomLogger.SuccessLogger.Info("User by ID:" + objectId.Hex() + " received successfully")
	return response, nil
}

func (handler *UserHandler) GetLoggedInUserInfo(ctx context.Context, request *pb.GetAllRequest) (*pb.User, error) {
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	user, err := handler.service.GetById(userId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + userId + " not found")
		return nil, err
	}
	pbUser := mapUser(user)
	handler.CustomLogger.SuccessLogger.Info("User received successfully")
	return pbUser, nil
}

func (handler *UserHandler) UpdateBasicInfo(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapBasicInfo(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Basic info updated successfully")
	return response, err
}

func (handler *UserHandler) UpdateExperienceAndEducation(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapExperienceAndEducation(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Experience and education updated successfully")
	return response, err
}

func (handler *UserHandler) UpdateSkillsAndInterests(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapSkillsAndInterests(mapUser(oldUser), request.User)
	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.CustomLogger.SuccessLogger.Info("Skills and interests updated successfully")
	return response, err
}

func (handler *UserHandler) GetEmail(ctx context.Context, request *pb.GetRequest) (*pb.GetEmailResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	user, err := handler.service.Get(objectId)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return nil, err
	}
	if !user.IsActive {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + request.Id + " is not activated")
		return nil, errors.New("Account is not activated")
	}
	handler.CustomLogger.SuccessLogger.Info("User email received successfully")
	response := &pb.GetEmailResponse{
		Email: user.Email,
	}
	return response, nil
}
func (handler *UserHandler) UpdateIsActiveById(ctx context.Context, request *pb.ActivateAccountRequest) (*pb.ActivateAccountResponse, error) {
	err := handler.service.UpdateIsActiveById(request.Id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + request.Id + " not activated")

		return &pb.ActivateAccountResponse{
			Success: err.Error(),
		}, err
	}
	handler.CustomLogger.SuccessLogger.Info("User by ID:" + request.Id + " received successfully")
	return &pb.ActivateAccountResponse{
		Success: "Success",
	}, nil
}

func (handler *UserHandler) GetIsActive(ctx context.Context, request *pb.GetRequest) (*pb.IsActiveResponse, error) {
	fmt.Println(request.Id)
	user, err := handler.service.GetById(request.Id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with ID:" + request.Id + " not found")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("User by id received successfully")
	return &pb.IsActiveResponse{
		IsActive: user.IsActive,
	}, nil
}

func (handler *UserHandler) GetIdByEmail(ctx context.Context, request *pb.GetIdByEmailRequest) (*pb.InsertResponse, error) {
	userId, err := handler.service.GetIdByEmail(request.Email)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with email:" + request.Email + " not found")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("User by email received successfully")
	return &pb.InsertResponse{
		Id: userId,
	}, nil
}

func (handler *UserHandler) GetIdByUsername(ctx context.Context, request *pb.GetIdByUsernameRequest) (*pb.InsertResponse, error) {
	user, err := handler.service.GetByUsername(request.Username)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("User with username:" + request.Username + " not found")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("User received successfully")
	return &pb.InsertResponse{
		Id: user.Id.Hex(),
	}, nil
}
