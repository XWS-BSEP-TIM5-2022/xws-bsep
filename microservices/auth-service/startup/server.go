package startup

import (
	"fmt"
	"log"
	"net"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/infrastructure/api"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/startup/config"
	"github.com/dgrijalva/jwt-go"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	auth_service_proto "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	saga "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/messaging"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/messaging/nats"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

const (
	QueueGroup = "auth_service"
)

func (server *Server) Start() {
	postgresClient := server.initPostgresClient()
	authStore := server.initAuthStore(postgresClient)

	jwtServiceClient, err := server.initJWTManager(server.config.PrivateKey, server.config.PublicKey)
	if err != nil {
		log.Fatal(err)
	}
	userServiceClient := server.initUserServiceClient()

	authService := server.initAuthService(authStore, userServiceClient, jwtServiceClient)

	commandSubscriber := server.initSubscriber(server.config.CreateUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.CreateUserReplySubject)
	server.initCreateUserHandler(authService, replyPublisher, commandSubscriber)

	authHandler := server.initAuthHandler(authService)

	server.startGrpcServer(authHandler)
}

func (server *Server) initPostgresClient() *gorm.DB {
	client, err := persistence.GetClient(
		server.config.AuthDBHost, server.config.AuthDBUser,
		server.config.AuthDBPass, server.config.AuthDBName,
		server.config.AuthDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initAuthStore(client *gorm.DB) *persistence.AuthPostgresStore {
	store, err := persistence.NewAuthPostgresStore(client)
	if err != nil {
		log.Fatal(err)
	}
	store.DeleteAll()
	for _, Auth := range auths {
		err := store.Insert(Auth)
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, Role := range roles {
		err := store.InsertRole(Role)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initPublisher(subject string) saga.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initSubscriber(subject, queueGroup string) saga.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}

func (server *Server) initAuthService(store *persistence.AuthPostgresStore, userServiceClient user.UserServiceClient, jwtService *application.JWTService) *application.AuthService {
	return application.NewAuthService(store, jwtService, userServiceClient)
}

func (server *Server) initCreateUserHandler(service *application.AuthService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewCreateUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initAuthHandler(service *application.AuthService) *api.AuthHandler {
	return api.NewAuthHandler(service)
}

func (server *Server) initJWTManager(privateKey, publicKey string) (*application.JWTService, error) {
	return application.NewJWTManager(privateKey, publicKey)
}

func (server *Server) initUserServiceClient() user.UserServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.UserServiceHost, server.config.UserServicePort)
	return persistence.NewUserServiceClient(address)
}

func (server *Server) startGrpcServer(authHandler *api.AuthHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// grpcServer := grpc.NewServer()
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		log.Fatalf("failed to parse public key: %v", err)
	}

	// interceptor := interceptor.NewAuthInterceptor(config.AccessibleRoles(), publicKey)
	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), publicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	auth_service_proto.RegisterAuthServiceServer(grpcServer, authHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
