package startup

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/infrastructure/api"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/startup/config"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	auth_service_proto "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	saga "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/messaging"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/messaging/nats"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/dgrijalva/jwt-go"
	otgo "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	config       *config.Config
	CustomLogger *api.CustomLogger
	tracer       otgo.Tracer
	closer       io.Closer
}

const name = "auth-service"

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

const (
	QueueGroup = "auth_service"
)

func (server *Server) Start() {
	postgresClient := server.initPostgresClient()
	authStore := server.initAuthStore(postgresClient)

	jwtServiceClient, err := server.initJWTManager(server.config.PrivateKey, server.config.PublicKey)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Initialization JWT service error")
		log.Fatal(err)
	}
	userServiceClient := server.initUserServiceClient()

	apiTokenServiceClient, err := server.initApiTokenManager(server.config.PrivateKeyApiToken, server.config.PublicKeyApiToken)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Initialization API service error")
		log.Fatal(err)
	}

	authService := server.initAuthService(authStore, userServiceClient, jwtServiceClient, apiTokenServiceClient)

	commandSubscriber := server.initSubscriber(server.config.CreateUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.CreateUserReplySubject)
	server.initCreateUserHandler(authService, replyPublisher, commandSubscriber)

	authHandler := server.initAuthHandler(authService)

	server.CustomLogger.SuccessLogger.Info("Starting auth service successfully, PORT: ", config.NewConfig().Port)
	server.startGrpcServer(authHandler)
}

func (server *Server) initPostgresClient() *gorm.DB {
	client, err := persistence.GetClient(
		server.config.AuthDBHost, server.config.AuthDBUser,
		server.config.AuthDBPass, server.config.AuthDBName,
		server.config.AuthDBPort)
	if err != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(server.config.AuthDBPass), bcrypt.DefaultCost)
		if err != nil {
			server.CustomLogger.ErrorLogger.Error("Starting the database failed because the password was not hashed")
		}
		server.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"auth_db_host":     server.config.AuthDBHost,
			"auth_db_port":     server.config.AuthDBPort,
			"auth_db_user":     server.config.AuthDBUser,
			"auth_db_password": string(hashedPassword), // TODO SD: password kao plain txt/hesirana?
			"auth_db_name":     server.config.AuthDBName,
		}).Error("Postgres database initialization error")
		// log.Fatal(err)
	}
	return client
}

func (server *Server) initAuthStore(client *gorm.DB) *persistence.AuthPostgresStore {
	store, err := persistence.NewAuthPostgresStore(client)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Auth store initialization error")
		log.Fatal(err)
	}
	store.DeleteAll()
	for _, Auth := range auths {
		err := store.Insert(Auth)
		if err != nil {
			server.CustomLogger.ErrorLogger.WithField("auth_id", Auth.Id).Error("Failed seed base with auth credentials")
			log.Fatal(err)
		}
	}
	for _, Role := range roles {
		err := store.InsertRole(Role)
		if err != nil {
			server.CustomLogger.ErrorLogger.WithField("role_id", Role.ID).Error("Failed seed base with user roles")
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initAuthService(store *persistence.AuthPostgresStore, userServiceClient user.UserServiceClient, jwtService *api.JWTService, apiTokenService *api.APITokenService) *api.AuthService {
	return api.NewAuthService(store, jwtService, userServiceClient, apiTokenService)
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

func (server *Server) initPublisher(subject string) saga.Publisher {
	log.Println(" *********** nats *************** ")
	log.Println(server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	log.Println(" ************************** ")
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initAuthHandler(service *api.AuthService) *application.AuthHandler {
	return application.NewAuthHandler(service)
}

func (server *Server) initJWTManager(privateKey, publicKey string) (*api.JWTService, error) {
	return api.NewJWTManager(privateKey, publicKey)
}
func (server *Server) initApiTokenManager(privateKey, publicKey string) (*api.APITokenService, error) {
	return api.NewAPITokenManager(privateKey, publicKey)
}

func (server *Server) initUserServiceClient() user.UserServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.UserServiceHost, server.config.UserServicePort)
	return persistence.NewUserServiceClient(address)
}

func (server *Server) startGrpcServer(authHandler *application.AuthHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to listen: %v", listener)
		// log.Fatalf("failed to listen: %v", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to parse public key")
		// log.Fatalf("failed to parse public key: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), publicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	auth_service_proto.RegisterAuthServiceServer(grpcServer, authHandler)
	if err := grpcServer.Serve(listener); err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to serve: %v", listener)
		// log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initCreateUserHandler(service *api.AuthService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewCreateUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}
