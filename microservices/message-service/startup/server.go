package startup

import (
	"fmt"
	connection "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	message "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/message_service"
	notification "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/sirupsen/logrus"
	"log"
	"net"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/message_service/startup/config"
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
	server.CustomLogger.SuccessLogger.Info("MongoDB initialization for message service successful, PORT: ", server.config.MessageDBPort, ", HOST: ", server.config.MessageDBHost)

	notificationServiceClient := server.initNotificationServiceClient()
	userServiceClient := server.initUserServiceClient()

	eventStore := server.initEventStore(mongoClient)

	messageStore := server.initMessageStore(mongoClient)
	messageService := server.initMessageService(messageStore, eventStore)
	messageHandler := server.initMessageHandler(messageService, notificationServiceClient, userServiceClient)

	server.CustomLogger.SuccessLogger.Info("Starting gRPC server for message service")

	server.startGrpcServer(messageHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.MessageDBHost, server.config.MessageDBPort)
	if err != nil {
		server.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"message_db_host": server.config.MessageDBHost,
			"message_db_port": server.config.MessageDBPort,
		}).Error("MongoDB initialization for message service failed")
	}
	return client
}

func (server *Server) initMessageStore(client *mongo.Client) domain.MessageStore {
	store := persistence.NewMessageMongoDBStore(client)
	//store.DeleteAll()
	//for _, message := range messages {
	//	_, err := store.Insert(message)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	return store
}

func (server *Server) initEventStore(client *mongo.Client) domain.EventStore {
	store := persistence.NewEventMongoDBStore(client)
	//store.DeleteAll()
	//for _, message := range messages {
	//	_, err := store.Insert(message)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	return store
}

func (server *Server) initMessageService(store domain.MessageStore, eventStore domain.EventStore) *application.MessageService {
	return application.NewMessageService(store, eventStore)
}

func (server *Server) initMessageHandler(service *application.MessageService, notificationServiceClient notification.NotificationServiceClient,
	userServiceClient user.UserServiceClient) *api.MessageHandler {
	return api.NewMessageHandler(service, notificationServiceClient, userServiceClient)
}

func (server *Server) initNotificationServiceClient() notification.NotificationServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.NotificationServiceHost, server.config.NotificationServicePort)
	return persistence.NewNotificationServiceClient(address)
}

func (server *Server) initUserServiceClient() user.UserServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.UserServiceHost, server.config.UserServicePort)
	return persistence.NewUserServiceClient(address)
}

func (server *Server) initConnectionServiceClient() connection.ConnectionServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.ConnectionServiceHost, server.config.ConnectionServicePort)
	return persistence.NewConnectionServiceClient(address)
}

func (server *Server) startGrpcServer(messageHandler *api.MessageHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to listen in message service: ", listener)
		log.Fatalf("failed to listen: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Parsing RSA public key in message service failed")
		log.Fatalf("failed to parse public key: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), publicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))

	message.RegisterMessageServiceServer(grpcServer, messageHandler)
	if err := grpcServer.Serve(listener); err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to serve gRPC in message service: ", listener)
		log.Fatalf("failed to serve: %s", err)
	}
}
