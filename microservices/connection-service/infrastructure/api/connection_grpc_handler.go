package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/application"
	"log"
	"regexp"
	"strconv"
)

type ConnectionHandler struct {
	pb.UnimplementedConnectionServiceServer
	service      *application.ConnectionService
	CustomLogger *CustomLogger
}

func NewConnectionHandler(service *application.ConnectionService) *ConnectionHandler {
	CustomLogger := NewCustomLogger()
	return &ConnectionHandler{
		service:      service,
		CustomLogger: CustomLogger,
	}
}

func (handler *ConnectionHandler) GetConnections(ctx context.Context, request *pb.GetRequest) (*pb.Users, error) {
	/* log injection prevention */
	id := request.UserID
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	id = re.ReplaceAllString(id, " ")

	//prosledili smo registrovanog korisnika
	//id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	friends, err := handler.service.GetConnections(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Connections for user with ID: " + id + " not found")
		return nil, err
	}
	response := &pb.Users{}
	for _, user := range friends {
		fmt.Println("User", id, "is friend with", user.UserID)
		response.Users = append(response.Users, mapUserConn(user))
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(friends)) + " connections for user with ID: " + id)
	return response, nil
}

func (handler *ConnectionHandler) GetRequests(ctx context.Context, request *pb.GetRequest) (*pb.Users, error) {
	/* log injection prevention */
	id := request.UserID
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	id = re.ReplaceAllString(id, " ")

	//prosledili smo registrovanog korisnika
	//id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	friends, err := handler.service.GetRequests(id)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Requests for user with ID: " + id + " not found")
		return nil, err
	}
	response := &pb.Users{}
	for _, user := range friends {
		fmt.Println("User", id, "has requests by", user.UserID)
		response.Users = append(response.Users, mapUserConn(user))
	}
	handler.CustomLogger.SuccessLogger.Info("Found " + strconv.Itoa(len(friends)) + " requests for user with ID: " + id)
	return response, nil
}

func (handler *ConnectionHandler) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.ActionResult, error) {
	userID := request.User.UserID
	isPublic := request.User.IsPublic

	/* log injection prevention */
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	userID = re.ReplaceAllString(userID, " ")

	register, err := handler.service.Register(userID, isPublic)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Registration for user with ID: " + userID + " failed")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("Registration for user with ID: " + userID + " successful")
	return register, err
}

func (handler *ConnectionHandler) AddConnection(ctx context.Context, request *pb.AddConnectionRequest) (*pb.AddConnectionResult, error) {
	//prosledili smo registrovanog korisnika
	userIDa := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	userIDb := request.AddConnectionDTO.UserID
	isPublic := request.AddConnectionDTO.IsPublic

	/* log injection prevention */
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	userIDb = re.ReplaceAllString(userIDb, " ")

	connection, err := handler.service.AddConnection(userIDa, userIDb, isPublic)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Creating connection between user with ID: " + userIDa + " and user with ID: " + userIDb + " failed")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("Creating connection between user with ID: " + userIDa + " and user with ID: " + userIDb + " successful")
	return connection, err
}

func (handler *ConnectionHandler) ApproveConnection(ctx context.Context, request *pb.ApproveConnectionRequest) (*pb.ActionResult, error) {
	//prosledili smo registrovanog korisnika
	userIDa := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	userIDb := request.ApproveConnectionDTO.UserID

	/* log injection prevention */
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	userIDb = re.ReplaceAllString(userIDb, " ")

	connection, err := handler.service.ApproveConnection(userIDa, userIDb)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Connection between user with ID: " + userIDa + " and user with ID: " + userIDb + " not approved")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("Approved connection between user with ID: " + userIDa + " and user with ID: " + userIDb)
	return connection, err
}

func (handler *ConnectionHandler) RejectConnection(ctx context.Context, request *pb.RejectConnectionRequest) (*pb.ActionResult, error) {
	//prosledili smo registrovanog korisnika
	userIDa := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	userIDb := request.RejectConnectionDTO.UserID

	/* log injection prevention */
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	userIDb = re.ReplaceAllString(userIDb, " ")

	connection, err := handler.service.RejectConnection(userIDa, userIDb)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Connection between user with ID: " + userIDa + " and user with ID: " + userIDb + " not rejected")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("Rejected connection between user with ID: " + userIDa + " and user with ID: " + userIDb)
	return connection, err
}

func (handler *ConnectionHandler) CheckConnection(ctx context.Context, request *pb.CheckConnectionRequest) (*pb.ConnectedResult, error) {
	//prosledili smo registrovanog korisnika
	userIDa := request.UserID
	userIDb := request.UserIDb

	/* log injection prevention */
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	userIDa = re.ReplaceAllString(userIDa, " ")
	userIDb = re.ReplaceAllString(userIDb, " ")

	connection, err := handler.service.CheckConnection(userIDa, userIDb)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Checking connection between user with ID: " + userIDa + " and user with ID: " + userIDb + " failed")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("Checking connection between user with ID: " + userIDa + " and user with ID: " + userIDb)
	return connection, err
}
