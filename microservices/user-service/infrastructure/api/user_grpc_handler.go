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
	"time"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/startup/config"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	service       *application.UserService
	WarningLogger *log.Logger
	InfoLogger    *log.Logger //*zerolog.Logger
	ErrorLogger   *log.Logger
	SuccessLogger *log.Logger
	DebugLogger   *log.Logger
	pb.UnimplementedUserServiceServer
}

func NewUserHandler(service *application.UserService) *UserHandler {
	InfoLogger := setLogger(config.NewConfig().InfoLogsFile, "INFO ")
	// InfoLogger := setZeroLogger(config.NewConfig().InfoLogsFile, "INFO ")
	ErrorLogger := setLogger(config.NewConfig().ErrorLogsFile, "ERROR ")
	WarningLogger := setLogger(config.NewConfig().WarningLogsFile, "WARNING ")
	SuccessLogger := setLogger(config.NewConfig().SuccessLogsFile, "SUCCESS ")
	DebugLogger := setLogger(config.NewConfig().DebugLogsFile, "DEBUG ")

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

func setZeroLogger(filename, loggerType string) *zerolog.Logger {
	logsFolderName := config.NewConfig().LogsFolder

	if _, err := os.Stat(logsFolderName); os.IsNotExist(err) {
		os.Mkdir(logsFolderName, 0777)
	}
	file, err := os.OpenFile(logsFolderName+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	// output := zerolog.ConsoleWriter{Out: mw, TimeFormat: time.RFC3339}
	// output.FormatLevel = func(i interface{}) string {
	// 	return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	// }
	// output.FormatMessage = func(i interface{}) string {
	// 	return fmt.Sprintf("%s", i)
	// }
	// output.FormatFieldName = func(i interface{}) string {
	// 	return fmt.Sprintf("%s:", i)
	// }
	// output.FormatFieldValue = func(i interface{}) string {
	// 	return strings.ToUpper(fmt.Sprintf("%s", i))
	// }

	output := zerolog.ConsoleWriter{Out: mw, TimeFormat: time.RFC3339Nano}
	Logger := zerolog.New(output).Output(output).With().Timestamp().Caller().Logger()
	return &Logger
}

func caller() func(*runtime.Frame) (function string, file string) {
	return func(f *runtime.Frame) (function string, file string) {
		p, _ := os.Getwd()

		return "", fmt.Sprintf("%s:%d", strings.TrimPrefix(f.File, p), f.Line)
	}
}

func setLogrusLogger() *logrus.Entry {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: caller(),
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyFile: "caller",
		},
	})
	logsFolderName := config.NewConfig().LogsFolder

	if _, err := os.Stat(logsFolderName); os.IsNotExist(err) {
		os.Mkdir(logsFolderName, 0777)
	}
	file, err := os.OpenFile(logsFolderName+"/info.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(mw)

	contextLogger := logrus.WithFields(logrus.Fields{})
	return contextLogger
}

func (handler *UserHandler) logMessage(msg string, loger *log.Logger) {
	// TODO SD: escape karaktera
	loger.Println(strings.ReplaceAll(msg, " ", "_"))
}

func (handler *UserHandler) GetAll(ctx context.Context, request *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	users, err := handler.service.GetAll()
	if err != nil {
		handler.logMessage("Get all", handler.ErrorLogger)
		return nil, err
	}
	response := &pb.GetAllResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}
	return response, nil
}

func (handler *UserHandler) GetAllPublic(ctx context.Context, request *pb.GetAllPublicRequest) (*pb.GetAllPublicResponse, error) {
	// TODO SD:
	// handler.InfoLogger.Info().Msgf(strings.ReplaceAll("Getting all public accounts", " ", "_"))
	// handler.InfoLogger.Error().Msgf(strings.ReplaceAll("Getting all public accounts test", " ", "_"))
	setLogrusLogger().WithFields(logrus.Fields{
		"CAO": "CAOCAO",
	}).Info("Aaaaaaaaa")

	users, err := handler.service.GetAllPublic()
	if err != nil {
		handler.logMessage("Found "+strconv.Itoa(len(users))+" public users", handler.ErrorLogger)
		return nil, err
	}
	response := &pb.GetAllPublicResponse{
		Users: []*pb.User{},
	}
	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}

	handler.logMessage("Found "+strconv.Itoa(len(users))+" public users", handler.SuccessLogger)
	return response, nil
}

func (handler *UserHandler) Search(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	criteria := request.Criteria
	users, err := handler.service.Search(criteria)

	if err != nil {
		handler.logMessage("Search error", handler.ErrorLogger)
		return nil, err
	}

	response := &pb.SearchResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		current := mapUser(user)
		response.Users = append(response.Users, current)
	}

	return response, nil
}
func (handler *UserHandler) Insert(ctx context.Context, request *pb.InsertRequest) (*pb.InsertResponse, error) {
	user := mapInsertUser(request.User)
	user, err := handler.service.Insert(user)

	if err != nil {
		handler.logMessage("User is not inserted", handler.ErrorLogger)
		return nil, err
	} else {
		return &pb.InsertResponse{
			Id: user.Id.Hex(),
		}, nil
	}
}

func (handler *UserHandler) Update(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	//id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(request.User.Id)
	if err != nil {
		handler.logMessage("ObjectId not created", handler.ErrorLogger)
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		handler.logMessage("User with ID:"+objectId.Hex()+" not found", handler.ErrorLogger)
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapUpdateUser(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	return response, err
}

func (handler *UserHandler) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.logMessage("ObjectId not created with ID:"+id, handler.ErrorLogger)
		return nil, err
	}
	user, err := handler.service.Get(objectId)
	if err != nil {
		handler.logMessage("User with ID:"+objectId.Hex()+" not found", handler.ErrorLogger)
		return nil, err
	}
	userPb := mapUser(user)
	response := &pb.GetResponse{
		User: userPb,
	}
	return response, nil
}

func (handler *UserHandler) GetLoggedInUserInfo(ctx context.Context, request *pb.GetAllRequest) (*pb.User, error) {
	userId := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	user, err := handler.service.GetById(userId)
	if err != nil {
		handler.logMessage("User with ID:"+userId+" not found", handler.ErrorLogger)
		return nil, err
	}
	pbUser := mapUser(user)
	return pbUser, nil
}

func (handler *UserHandler) UpdateBasicInfo(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.logMessage("ObjectId not created with ID:"+id, handler.ErrorLogger)
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		handler.logMessage("User with ID:"+objectId.Hex()+" not found", handler.ErrorLogger)
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapBasicInfo(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	return response, err
}

func (handler *UserHandler) UpdateExperienceAndEducation(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.logMessage("ObjectId not created with ID:"+id, handler.ErrorLogger)
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		handler.logMessage("User with ID:"+objectId.Hex()+" not found", handler.ErrorLogger)
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapExperienceAndEducation(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	return response, err
}

func (handler *UserHandler) UpdateSkillsAndInterests(ctx context.Context, request *pb.UpdateRequest) (*pb.UpdateResponse, error) {

	id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.logMessage("ObjectId not created with ID:"+id, handler.ErrorLogger)
		return nil, err
	}
	oldUser, err := handler.service.Get(objectId)

	if err != nil {
		handler.logMessage("User with ID:"+objectId.Hex()+" not found", handler.ErrorLogger)
		return &pb.UpdateResponse{
			Success: "error",
		}, err
	}

	user := mapSkillsAndInterests(mapUser(oldUser), request.User)

	success, err := handler.service.Update(user)
	response := &pb.UpdateResponse{
		Success: success,
	}
	return response, err
}

func (handler *UserHandler) GetEmail(ctx context.Context, request *pb.GetRequest) (*pb.GetEmailResponse, error) {
	id := request.Id
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		handler.logMessage("ObjectId not created with ID:"+id, handler.ErrorLogger)
		return nil, err
	}
	user, err := handler.service.Get(objectId)
	if err != nil {
		handler.logMessage("User with ID:"+objectId.Hex()+" not found", handler.ErrorLogger)
		return nil, err
	}

	if !user.IsActive {
		handler.logMessage("User with ID:"+request.Id+" is not activated", handler.ErrorLogger)
		return nil, errors.New("Account is not activated")
	}

	response := &pb.GetEmailResponse{
		Email: user.Email,
	}
	return response, nil
}
func (handler *UserHandler) UpdateIsActiveById(ctx context.Context, request *pb.ActivateAccountRequest) (*pb.ActivateAccountResponse, error) {
	err := handler.service.UpdateIsActiveById(request.Id)
	if err != nil {
		handler.logMessage("User with ID:"+request.Id+" not activated", handler.ErrorLogger)

		return &pb.ActivateAccountResponse{
			Success: err.Error(),
		}, err
	}
	return &pb.ActivateAccountResponse{
		Success: "Success",
	}, nil
}

func (handler *UserHandler) GetIsActive(ctx context.Context, request *pb.GetRequest) (*pb.IsActiveResponse, error) {
	fmt.Println(request.Id)
	user, err := handler.service.GetById(request.Id)
	if err != nil {
		handler.logMessage("User with ID:"+request.Id+" not found", handler.ErrorLogger)
		return nil, err
	}
	return &pb.IsActiveResponse{
		IsActive: user.IsActive,
	}, nil
}

func (handler *UserHandler) GetIdByEmail(ctx context.Context, request *pb.GetIdByEmailRequest) (*pb.InsertResponse, error) {
	userId, err := handler.service.GetIdByEmail(request.Email)
	if err != nil {
		handler.logMessage("User with email:"+request.Email+" not found", handler.ErrorLogger)
		return nil, err
	}
	return &pb.InsertResponse{
		Id: userId,
	}, nil
}

func (handler *UserHandler) GetIdByUsername(ctx context.Context, request *pb.GetIdByUsernameRequest) (*pb.InsertResponse, error) {
	user, err := handler.service.GetByUsername(request.Username)
	if err != nil {
		handler.logMessage("User with username:"+request.Username+" not found", handler.ErrorLogger)
		return nil, err
	}
	return &pb.InsertResponse{
		Id: user.Id.Hex(),
	}, nil
}
