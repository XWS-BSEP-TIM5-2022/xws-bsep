package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/startup/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	service       *application.UserService
	WarningLogger *logrus.Entry
	InfoLogger    *logrus.Entry //*log.Logger //*zerolog.Logger
	ErrorLogger   *logrus.Entry
	SuccessLogger *logrus.Entry
	DebugLogger   *logrus.Entry
	pb.UnimplementedUserServiceServer
}

func NewUserHandler(service *application.UserService) *UserHandler {
	InfoLogger := setLogrusLogger(config.NewConfig().InfoLogsFile)
	ErrorLogger := setLogrusLogger(config.NewConfig().ErrorLogsFile)
	WarningLogger := setLogrusLogger(config.NewConfig().WarningLogsFile)
	SuccessLogger := setLogrusLogger(config.NewConfig().SuccessLogsFile)
	DebugLogger := setLogrusLogger(config.NewConfig().DebugLogsFile)

	return &UserHandler{
		service:       service,
		InfoLogger:    InfoLogger,
		ErrorLogger:   ErrorLogger,
		WarningLogger: WarningLogger,
		SuccessLogger: SuccessLogger,
		DebugLogger:   DebugLogger,
	}
}

func setLogger(filename, loggerType string) *log.Logger {
	logsFolderName := config.NewConfig().LogsFolder
	if _, err := os.Stat(logsFolderName); os.IsNotExist(err) {
		os.Mkdir(logsFolderName, 0777)
	}
	file, err := os.OpenFile(logsFolderName+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	Logger := log.New(file, loggerType, log.Ldate|log.Ltime|log.Lshortfile) //Llongfile

	mw := io.MultiWriter(os.Stdout, file)
	Logger.SetOutput(mw)
	return Logger
}

func caller() func(*runtime.Frame) (function string, file string) {
	return func(f *runtime.Frame) (function string, file string) {
		p, _ := os.Getwd()

		return "", fmt.Sprintf("%s:%d", strings.TrimPrefix(f.File, p), f.Line)
	}
}

func setLogrusLogger(filename string) *logrus.Entry {
	mLog := logrus.New()
	mLog.SetReportCaller(true)

	logsFolderName := config.NewConfig().LogsFolder

	if _, err := os.Stat(logsFolderName); os.IsNotExist(err) {
		os.Mkdir(logsFolderName, 0777)
	}
	file, err := os.OpenFile(logsFolderName+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	mLog.SetOutput(mw)

	mLog.SetFormatter(&logrus.JSONFormatter{ //TextFormatter //JSONFormatter
		CallerPrettyfier: caller(),
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyFile: "mehtod",
		},
		// ForceColors: true,
	})
	contextLogger := mLog.WithFields(logrus.Fields{})
	return contextLogger
}

func (handler *UserHandler) logMessage(msg string, loger *log.Logger) {
	// TODO SD: escape karaktera
	loger.Println(strings.ReplaceAll(msg, " ", "_"))
}

func (handler *UserHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	users, err := handler.service.GetAll()
	if err != nil {
		handler.ErrorLogger.Error("Get all")
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	handler.SuccessLogger.Info("Found " + strconv.Itoa(len(users)) + " public users")
	return response, nil
}

func (handler *UserHandler) GetAllPublic(ctx context.Context, request *pb.GetAllPublicRequest) (*pb.GetAllPublicResponse, error) {
	handler.InfoLogger.WithFields(logrus.Fields{}).Info("Getting all public accounts")

	users, err := handler.service.GetAllPublic()
	if err != nil {
		handler.ErrorLogger.Error("Found " + strconv.Itoa(len(users)) + " public users")
		return nil, err
	}
	response := &pb.GetAllPublicResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	handler.SuccessLogger.Info("Found " + strconv.Itoa(len(users)) + " public users")
	return response, nil
}

func (handler *UserHandler) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	criteria := request.Criteria
	users, err := handler.service.Search(criteria)

	if err != nil {
		handler.ErrorLogger.Error("Search error")
		return nil, err
	}

	response := &pb.SearchResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	handler.SuccessLogger.Info("Number of users found after search: " + strconv.Itoa(len(users)))
	return response, nil
}
func (handler *UserHandler) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResponse, error) {
	user := mapInsertUser(request.User)
	user, err := handler.service.Insert(user)

	if err != nil {
		handler.ErrorLogger.Error("User is not inserted")
		return nil, err
	} else {
		handler.SuccessLogger.Info("User inserted successfully")
		return &pb.InsertResponse{
			Id: user.Id.Hex(),
		}, nil
	}
}

func (handler *UserHandler) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	//id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(request.User.Id)
	if err != nil {
		handler.ErrorLogger.Error("ObjectId not created")
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		handler.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapUpdateUser(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.SuccessLogger.Info("User updated successfully")
	return response, err
}

func (handler *UserHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	user, err := handler.service.Get(objectId)
	if err != nil {
		handler.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return nil, err
	}
	userPb := mapUser(user)
	response := &pb.GetResponse{
		User: userPb,
	}
	handler.SuccessLogger.Info("User by ID:" + objectId.Hex() + " received successfully")

	return response, nil
}

func (handler *UserHandler) GetLoggedInUserInfo(ctx context.Context, request *pb.GetAllRequest) (*pb.User, error) {
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	user, err := handler.service.GetById(userId)
	if err != nil {
		handler.ErrorLogger.Error("User with ID:" + userId + " not found")
		return nil, err
	}
	pbUser := mapUser(user)
	handler.SuccessLogger.Info("User received successfully")
	return pbUser, nil
}

func (handler *UserHandler) UpdateBasicInfo(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		handler.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapBasicInfo(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.SuccessLogger.Info("Basic info updated successfully")
	return response, err
}

func (handler *UserHandler) UpdateExperienceAndEducation(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		handler.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapExperienceAndEducation(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.SuccessLogger.Info("Experience and education updated successfully")
	return response, err
}

func (handler *UserHandler) UpdateSkillsAndInterests(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		handler.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapSkillsAndInterests(mapUser(oldUser), request.User)
	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	handler.SuccessLogger.Info("Skills and interests updated successfully")
	return response, err
}

func (handler *UserHandler) GetEmail(ctx context.Context, request *pb.GetRequest) (*pb.GetEmailResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.ErrorLogger.Error("ObjectId not created with ID:" + id)
		return nil, err
	}
	user, err := handler.service.Get(objectId)
	if err != nil {
		handler.ErrorLogger.Error("User with ID:" + objectId.Hex() + " not found")
		return nil, err
	}
	if !user.IsActive {
		handler.ErrorLogger.Error("User with ID:" + request.Id + " is not activated")
		return nil, errors.New("Account is not activated")
	}
	handler.SuccessLogger.Info("User email received successfully")
	response := &pb.GetEmailResponse{
		Email: user.Email,
	}
	return response, nil
}
func (handler *UserHandler) UpdateIsActiveById(ctx context.Context, request *pb.ActivateAccountRequest) (*pb.ActivateAccountResponse, error) {
	err := handler.service.UpdateIsActiveById(request.Id)
	if err != nil {
		handler.ErrorLogger.Error("User with ID:" + request.Id + " not activated")

		return &pb.ActivateAccountResponse{
			Success: err.Error(),
		}, err
	}
	handler.SuccessLogger.Info("User by ID:" + request.Id + " received successfully")
	return &pb.ActivateAccountResponse{
		Success: "Success",
	}, nil
}

func (handler *UserHandler) GetIsActive(ctx context.Context, request *pb.GetRequest) (*pb.IsActiveResponse, error) {
	fmt.Println(request.Id)
	user, err := handler.service.GetById(request.Id)
	if err != nil {
		handler.ErrorLogger.Error("User with ID:" + request.Id + " not found")
		return nil, err
	}
	handler.SuccessLogger.Info("User by id received successfully")
	return &pb.IsActiveResponse{
		IsActive: user.IsActive,
	}, nil
}

func (handler *UserHandler) GetIdByEmail(ctx context.Context, request *pb.GetIdByEmailRequest) (*pb.InsertResponse, error) {
	userId, err := handler.service.GetIdByEmail(request.Email)
	if err != nil {
		handler.ErrorLogger.Error("User with email:" + request.Email + " not found")
		return nil, err
	}
	handler.SuccessLogger.Info("User by email received successfully")
	return &pb.InsertResponse{
		Id: userId,
	}, nil
}

func (handler *UserHandler) GetIdByUsername(ctx context.Context, request *pb.GetIdByUsernameRequest) (*pb.InsertResponse, error) {
	user, err := handler.service.GetByUsername(request.Username)
	if err != nil {
		handler.ErrorLogger.Error("User with username:" + request.Username + " not found")
		return nil, err
	}
	handler.SuccessLogger.Info("User received successfully")
	return &pb.InsertResponse{
		Id: user.Id.Hex(),
	}, nil
}
