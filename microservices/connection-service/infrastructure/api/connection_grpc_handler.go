package api

import (
	"context"
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	notification "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/application"
	"log"
	"regexp"
	"strconv"
)

type ConnectionHandler struct {
	pb.UnimplementedConnectionServiceServer
	service                   *application.ConnectionService
	CustomLogger              *CustomLogger
	notificationServiceClient notification.NotificationServiceClient
	userServiceClient         user.UserServiceClient
}

func NewConnectionHandler(service *application.ConnectionService, notificationServiceClient notification.NotificationServiceClient,
	userServiceClient user.UserServiceClient) *ConnectionHandler {
	CustomLogger := NewCustomLogger()
	return &ConnectionHandler{
		service:                   service,
		CustomLogger:              CustomLogger,
		notificationServiceClient: notificationServiceClient,
		userServiceClient:         userServiceClient,
	}
}

func (handler *ConnectionHandler) GetConnections(ctx context.Context, request *pb.GetRequest) (*pb.Users, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConnections")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	/* log injection prevention */
	id := request.UserID
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	id = re.ReplaceAllString(id, " ")

	//prosledili smo registrovanog korisnika
	//id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	friends, err := handler.service.GetConnections(ctx, id)
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
	span := tracer.StartSpanFromContext(ctx, "GetRequests")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(ctx, span)

	/* log injection prevention */
	id := request.UserID
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	id = re.ReplaceAllString(id, " ")

	//prosledili smo registrovanog korisnika
	//id := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	friends, err := handler.service.GetRequests(ctx, id)
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
	span := tracer.StartSpanFromContext(ctx, "Register")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(ctx, span)

	userID := request.User.UserID
	isPublic := request.User.IsPublic

	/* log injection prevention */
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	userID = re.ReplaceAllString(userID, " ")

	register, err := handler.service.Register(ctx, userID, isPublic)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Registration for user with ID: " + userID + " failed")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("Registration for user with ID: " + userID + " successful")
	return register, err
}

func (handler *ConnectionHandler) AddConnection(ctx context.Context, request *pb.AddConnectionRequest) (*pb.AddConnectionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "AddConnection")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(ctx, span)

	//prosledili smo registrovanog korisnika
	userIDa := ctx.Value(interceptor.LoggedInUserKey{}).(string) // onaj koji sanje zahtev
	userIDb := request.AddConnectionDTO.UserID                   // onaj koji prima zahtev (njemu se salje notifikacija)
	isPublic := request.AddConnectionDTO.IsPublic
	isPublicLogged := request.AddConnectionDTO.IsPublicLogged

	/* log injection prevention */
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	userIDb = re.ReplaceAllString(userIDb, " ")

	connection, err := handler.service.AddConnection(ctx, userIDa, userIDb, isPublic, isPublicLogged)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Creating connection between user with ID: " + userIDa + " and user with ID: " + userIDb + " failed")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("Creating connection between user with ID: " + userIDa + " and user with ID: " + userIDb + " successful")

	// slanje notifikacija
	sender, _ := handler.userServiceClient.Get(ctx, &user.GetRequest{Id: userIDa})
	reciever, _ := handler.userServiceClient.Get(ctx, &user.GetRequest{Id: userIDb})
	if sender.User.PostNotification == true && reciever.User.IsPublic == false {
		notificationRequest := &notification.InsertNotificationRequest{}
		notificationRequest.Notification = &notification.Notification{}
		notificationRequest.Notification.Type = notification.Notification_NotificationTypeEnum(1)
		notificationRequest.Notification.Text = "User " + sender.User.Name + " " + sender.User.LastName + " requested to follow you"
		notificationRequest.Notification.UserId = userIDb
		handler.notificationServiceClient.Insert(ctx, notificationRequest)
	} else if sender.User.PostNotification == true && reciever.User.IsPublic == true {
		notificationRequest := &notification.InsertNotificationRequest{}
		notificationRequest.Notification = &notification.Notification{}
		notificationRequest.Notification.Type = notification.Notification_NotificationTypeEnum(1)
		notificationRequest.Notification.Text = "User " + sender.User.Name + " " + sender.User.LastName + " started following you"
		notificationRequest.Notification.UserId = userIDb
		handler.notificationServiceClient.Insert(ctx, notificationRequest)
	}

	return connection, err
}

func (handler *ConnectionHandler) BlockUser(ctx context.Context, request *pb.BlockUserRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "BlockUser")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(ctx, span)

	fmt.Println("BlockUser")
	//prosledili smo registrovanog korisnika
	userIDa := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	userIDb := request.BlockUserDTO.UserID
	return handler.service.BlockUser(ctx, userIDa, userIDb, request.BlockUserDTO.IsPublic, request.BlockUserDTO.IsPublicLogged)
}

func (handler *ConnectionHandler) ApproveConnection(ctx context.Context, request *pb.ApproveConnectionRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "ApproveConnection")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(ctx, span)

	//prosledili smo registrovanog korisnika
	userIDa := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	userIDb := request.ApproveConnectionDTO.UserID

	/* log injection prevention */
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	userIDb = re.ReplaceAllString(userIDb, " ")

	connection, err := handler.service.ApproveConnection(ctx, userIDa, userIDb)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Connection between user with ID: " + userIDa + " and user with ID: " + userIDb + " not approved")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("Approved connection between user with ID: " + userIDa + " and user with ID: " + userIDb)

	// slanje notifikacija
	sender, _ := handler.userServiceClient.Get(ctx, &user.GetRequest{Id: userIDa})
	if sender.User.PostNotification == true {
		notificationRequest := &notification.InsertNotificationRequest{}
		notificationRequest.Notification = &notification.Notification{}
		notificationRequest.Notification.Type = notification.Notification_NotificationTypeEnum(1)
		notificationRequest.Notification.Text = "User " + sender.User.Name + " " + sender.User.LastName + " accepted your follow request"
		notificationRequest.Notification.UserId = userIDb
		handler.notificationServiceClient.Insert(ctx, notificationRequest)
	}

	return connection, err
}

func (handler *ConnectionHandler) RejectConnection(ctx context.Context, request *pb.RejectConnectionRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "RejectConnection")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(ctx, span)

	//prosledili smo registrovanog korisnika
	userIDa := ctx.Value(interceptor.LoggedInUserKey{}).(string)
	userIDb := request.RejectConnectionDTO.UserID

	/* log injection prevention */
	re, err := regexp.Compile(`[^\w]`) // specijalni karakteri
	if err != nil {
		log.Fatal(err)
	}
	userIDb = re.ReplaceAllString(userIDb, " ")

	connection, err := handler.service.RejectConnection(ctx, userIDa, userIDb)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Connection between user with ID: " + userIDa + " and user with ID: " + userIDb + " not rejected")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("Rejected connection between user with ID: " + userIDa + " and user with ID: " + userIDb)
	return connection, err
}

func (handler *ConnectionHandler) CheckConnection(ctx context.Context, request *pb.CheckConnectionRequest) (*pb.ConnectedResult, error) {
	span := tracer.StartSpanFromContext(ctx, "CheckConnection")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

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

	connection, err := handler.service.CheckConnection(ctx, userIDa, userIDb)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Checking connection between user with ID: " + userIDa + " and user with ID: " + userIDb + " failed")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("Checking connection between user with ID: " + userIDa + " and user with ID: " + userIDb)
	return connection, err
}

func (handler *ConnectionHandler) GetRecommendation(ctx context.Context, request *pb.GetRequest) (*pb.Users, error) {
	span := tracer.StartSpanFromContext(ctx, "GetRecommendation")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(context.Background(), span)

	fmt.Println("GetRecommendation")

	id := request.UserID
	fmt.Println(id)
	recommendation, err := handler.service.GetRecommendation(ctx, id)
	if err != nil {
		return nil, err
	}
	response := &pb.Users{}
	for _, user := range recommendation {
		response.Users = append(response.Users, mapUserConn(user))
	}
	return response, nil
}

func (handler *ConnectionHandler) ChangePrivacy(ctx context.Context, request *pb.ChangePrivacyRequest) (*pb.ActionResult, error) {
	span := tracer.StartSpanFromContext(ctx, "ChangePrivacy")
	defer span.Finish()

	ctx = tracer.ContextWithSpan(ctx, span)

	//prosledili smo registrovanog korisnika
	userIDa := ctx.Value(interceptor.LoggedInUserKey{}).(string)

	connection, err := handler.service.ChangePrivacy(ctx, userIDa, request.ChangePrivacyDTO.Private)
	if err != nil {
		handler.CustomLogger.ErrorLogger.Error("Privacy of user with ID: " + userIDa + " successfully changed!")
		return nil, err
	}
	handler.CustomLogger.SuccessLogger.Info("Error while changing privacy of user with ID: " + userIDa)
	return connection, err
}
