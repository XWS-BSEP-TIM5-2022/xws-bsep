package startup

import (
	"fmt"
	connection "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	notification "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/sirupsen/logrus"
	"log"
	"net"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	post "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/startup/config"
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
	server.CustomLogger.SuccessLogger.Info("MongoDB initialization for post service successful, PORT: ", server.config.PostDBPort, ", HOST: ", server.config.PostDBHost)

	notificationServiceClient := server.initNotificationServiceClient()
	connectionServiceClient := server.initConnectionServiceClient()
	userServiceClient := server.initUserServiceClient()

	eventStore := server.initEventStore(mongoClient)
	postStore := server.initPostStore(mongoClient)
	postService := server.initPostService(postStore, eventStore)
	postHandler := server.initPostHandler(postService, notificationServiceClient, connectionServiceClient, userServiceClient)

	server.CustomLogger.SuccessLogger.Info("Starting gRPC server for post service")
	server.startGrpcServer(postHandler)
}

func (server *Server) initMongoClient() *mongo.Client { // inicijalizacija mongo klijenta
	client, err := persistence.GetClient(server.config.PostDBHost, server.config.PostDBPort)
	if err != nil {
		server.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"post_db_host": server.config.PostDBHost,
			"post_db_port": server.config.PostDBPort,
		}).Error("MongoDB initialization for post service failed")
	}
	return client
}

func (server *Server) initPostStore(client *mongo.Client) domain.PostStore { // inicijalizacija mongo baze
	store := persistence.NewPostMongoDBStore(client)
	//store.DeleteAll()
	//for _, post := range posts {
	//	_, err := store.Insert(post)
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

func (server *Server) initPostService(store domain.PostStore, eventStore domain.EventStore) *application.PostService {
	return application.NewPostService(store, eventStore)
}

func (server *Server) initNotificationServiceClient() notification.NotificationServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.NotificationServiceHost, server.config.NotificationServicePort)
	return persistence.NewNotificationServiceClient(address)
}

func (server *Server) initConnectionServiceClient() connection.ConnectionServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.ConnectionServiceHost, server.config.ConnectionServicePort)
	return persistence.NewConnectionServiceClient(address)
}

func (server *Server) initUserServiceClient() user.UserServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.UserServiceHost, server.config.UserServicePort)
	return persistence.NewUserServiceClient(address)
}

func (server *Server) initPostHandler(service *application.PostService, notificationServiceClient notification.NotificationServiceClient,
	connectionServiceClient connection.ConnectionServiceClient, userServiceClient user.UserServiceClient) *api.PostHandler {
	return api.NewPostHandler(service, notificationServiceClient, connectionServiceClient, userServiceClient)
}

func (server *Server) startGrpcServer(postHandler *api.PostHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to listen in post service: ", listener)
		log.Fatalf("failed to listen: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Parsing RSA public key in post service failed")
		log.Fatalf("failed to parse public key: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), publicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))

	post.RegisterPostServiceServer(grpcServer, postHandler)
	if err := grpcServer.Serve(listener); err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to serve gRPC in post service: ", listener)
		log.Fatalf("failed to serve: %s", err)
	}
}
