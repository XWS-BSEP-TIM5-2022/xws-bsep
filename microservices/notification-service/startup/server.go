package startup

import (
	"fmt"
	notification "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	"github.com/sirupsen/logrus"
	"log"
	"net"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/notification_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/notification_service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/notification_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/notification_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/notification_service/startup/config"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type Server struct {
	config       *config.Config
	CustomLogger *api.CustomLogger
}

func NewServer(config *config.Config) *Server {
	CustomLogger := api.NewCustomLogger()
	return &Server{
		config:       config,
		CustomLogger: CustomLogger,
	}
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	server.CustomLogger.SuccessLogger.Info("MongoDB initialization for notification service successful, PORT: ", server.config.NotificationDBPort, ", HOST: ", server.config.NotificationDBHost)

	notificationStore := server.initNotificationStore(mongoClient)
	notificationService := server.initNotificationService(notificationStore)
	notificationHandler := server.initNotificationHandler(notificationService)

	server.CustomLogger.SuccessLogger.Info("Starting gRPC server for notification service")
	server.startGrpcServer(notificationHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.NotificationDBHost, server.config.NotificationDBPort)
	if err != nil {
		server.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"notification_db_host": server.config.NotificationDBHost,
			"notification_db_port": server.config.NotificationDBPort,
		}).Error("MongoDB initialization for notification service failed")
	}
	return client
}

func (server *Server) initNotificationStore(client *mongo.Client) domain.NotificationStore {
	store := persistence.NewNotificationMongoDBStore(client)
	//store.DeleteAll()
	//for _, notification := range notifications {
	//	_, err := store.Insert(notification)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	return store
}

func (server *Server) initNotificationService(store domain.NotificationStore) *application.NotificationService {
	return application.NewNotificationService(store)
}

func (server *Server) initNotificationHandler(service *application.NotificationService) *api.NotificationHandler {
	return api.NewNotificationHandler(service)
}

func (server *Server) startGrpcServer(notificationHandler *api.NotificationHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to listen in notification service: ", listener)
		log.Fatalf("failed to listen: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Parsing RSA public key in notification service failed")
		log.Fatalf("failed to parse public key: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), publicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))

	notification.RegisterNotificationServiceServer(grpcServer, notificationHandler)
	if err := grpcServer.Serve(listener); err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to serve gRPC in notification service: ", listener)
		log.Fatalf("failed to serve: %s", err)
	}
}
