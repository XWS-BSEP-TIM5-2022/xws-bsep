package startup

import (
	"fmt"
	connection "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	message "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/message_service"
	notification "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	post "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	otgo "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"io"
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
	tracer       otgo.Tracer
	closer       io.Closer
}

const name = "notification-service"

func NewServer(config *config.Config) *Server {
	CustomLogger := api.NewCustomLogger()
	tracer, closer := tracer.Init(name)
	otgo.SetGlobalTracer(tracer)

	return &Server{
		config:       config,
		CustomLogger: CustomLogger,
		tracer:       tracer,
		closer:       closer,
	}
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	server.CustomLogger.SuccessLogger.Info("MongoDB initialization for notification service successful, PORT: ", server.config.NotificationDBPort, ", HOST: ", server.config.NotificationDBHost)

	eventStore := server.initEventStore(mongoClient)
	notificationStore := server.initNotificationStore(mongoClient)
	notificationService := server.initNotificationService(notificationStore, eventStore)
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

func (server *Server) initNotificationService(store domain.NotificationStore, eventStore domain.EventStore) *application.NotificationService {
	return application.NewNotificationService(store, eventStore)
}

func (server *Server) initUserServiceClient() user.UserServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.UserServiceHost, server.config.UserServicePort)
	return persistence.NewUserServiceClient(address)
}

func (server *Server) initPostServiceClient() post.PostServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.PostServiceHost, server.config.PostServicePort)
	return persistence.NewPostServiceClient(address)
}

func (server *Server) initConnectionServiceClient() connection.ConnectionServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.ConnectionServiceHost, server.config.ConnectionServicePort)
	return persistence.NewConnectionServiceClient(address)
}

func (server *Server) initMessageServiceClient() message.MessageServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.MessageServiceHost, server.config.MessageServicePort)
	return persistence.NewMessageServiceClient(address)
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
