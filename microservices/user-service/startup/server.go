package startup

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/infrastructure/persistence"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	interceptor "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	inventory "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	saga "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/messaging"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/saga/messaging/nats"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/startup/config"
	otgo "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type Server struct {
	config       *config.Config
	CustomLogger *api.CustomLogger
	tracer       otgo.Tracer
	closer       io.Closer
}

const name = "user-service"

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
	QueueGroup = "user_service"
)

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	userStore := server.initUserStore(mongoClient)

	commandPublisher := server.initPublisher(server.config.CreateUserCommandSubject)
	replySubscriber := server.initSubscriber(server.config.CreateUserReplySubject, QueueGroup)
	createUserOrchestrator := server.initCreateUserOrchestrator(commandPublisher, replySubscriber)

	userService := server.initUserService(userStore, createUserOrchestrator)

	commandSubscriber := server.initSubscriber(server.config.CreateUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.CreateUserReplySubject)
	server.initCreateUserHandler(userService, replyPublisher, commandSubscriber)

	userHandler := server.initUserHandler(userService)

	server.CustomLogger.SuccessLogger.Info("Starting user service successfully, PORT: ", config.NewConfig().Port)
	server.startGrpcServer(userHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.UserDBHost, server.config.UserDBPort)
	if err != nil {
		server.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"user_db_host": server.config.UserDBHost,
			"user_db_port": server.config.UserDBPort,
		}).Error("Mongo database initialization error")
		// log.Fatal(err)
	}
	return client
}

func (server *Server) initUserStore(client *mongo.Client) domain.UserStore {
	store := persistence.NewUserMongoDBStore(client)
	store.DeleteAll()
	for _, user := range users {
		_, err := store.Insert(user)
		if err != nil {
			server.CustomLogger.ErrorLogger.Error("User store initialization error")
			// log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initPublisher(subject string) saga.Publisher {
	log.Println(server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
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

func (server *Server) initCreateUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) *application.CreateUserOrchestrator {
	orchestrator, err := application.NewCreateUserOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) initUserService(store domain.UserStore, orchestrator *application.CreateUserOrchestrator) *application.UserService {
	return application.NewUserService(store, orchestrator)
}

func (server *Server) initCreateUserHandler(service *application.UserService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := api.NewCreateUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initUserHandler(service *application.UserService) *api.UserHandler {
	return api.NewUserHandler(service)
}

func (server *Server) startGrpcServer(userHandler *api.UserHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to listen: %v", listener)
		log.Fatalf("failed to listen: %v", err)
	}
	// ****
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to parse public key")
		// log.Fatalf("failed to parse public key: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), publicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	// ***
	inventory.RegisterUserServiceServer(grpcServer, userHandler)
	if err := grpcServer.Serve(listener); err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to serve: %v", listener)
		// log.Fatalf("failed to serve: %s", err)
	}
}
