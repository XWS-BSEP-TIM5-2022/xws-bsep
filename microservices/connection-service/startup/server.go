package startup

import (
	"fmt"
	notification "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"io"
	"log"
	"net"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	"github.com/dgrijalva/jwt-go"
	otgo "github.com/opentracing/opentracing-go"

	inventory "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/startup/config"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"google.golang.org/grpc"
)

type Server struct {
	config       *config.Config
	CustomLogger *api.CustomLogger
	tracer       otgo.Tracer
	closer       io.Closer
}

const name = "connection-service"

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
	QueueGroup = "connection_service"
)

func (server *Server) Start() {
	neo4jClient := server.initNeo4J()
	server.CustomLogger.SuccessLogger.Info("Neo4J initialization for connection service successful, PORT: ", server.config.Port)

	notificationServiceClient := server.initNotificationServiceClient()
	userServiceClient := server.initUserServiceClient()

	connectionStore := server.initConnectionStore(neo4jClient)
	connectionService := server.initConnectionService(connectionStore)
	connectionHandler := server.initConnectionHandler(connectionService, notificationServiceClient, userServiceClient)

	server.CustomLogger.SuccessLogger.Info("Starting gRPC server for connection service")
	server.startGrpcServer(connectionHandler)
}

func (server *Server) initNeo4J() *neo4j.Driver {

	dbUri := "bolt://" + server.config.ConnectionDBHost + ":" + server.config.ConnectionDBPort
	// dbUri := "bolt://localhost:7687"
	server.CustomLogger.InfoLogger.Info("Neo4J datase on " + dbUri)

	client, err := persistence.GetClient(dbUri, server.config.ConnectionDBUser, server.config.ConnectionDBPass)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Neo4J initialization for connection service failed")
		log.Fatal(err)
	}
	return client
}

func (server *Server) initConnectionStore(client *neo4j.Driver) domain.ConnectionStore {
	store := persistence.NewConnectionDBStore(client)
	/*
		store.DeleteAll()
		for _, product := range products {
			err := store.Insert(product)
			if err != nil {
				log.Fatal(err)
			}
		}
	*/
	return store
}

func (server *Server) initConnectionService(store domain.ConnectionStore) *application.ConnectionService {
	return application.NewConnectionService(store)
}

func (server *Server) initConnectionHandler(service *application.ConnectionService, notificationServiceClient notification.NotificationServiceClient,
	userServiceClient user.UserServiceClient) *api.ConnectionHandler {
	return api.NewConnectionHandler(service, notificationServiceClient, userServiceClient)
}

func (server *Server) initNotificationServiceClient() notification.NotificationServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.NotificationServiceHost, server.config.NotificationServicePort)
	return persistence.NewNotificationServiceClient(address)
}

func (server *Server) initUserServiceClient() user.UserServiceClient {
	address := fmt.Sprintf("%s:%s", server.config.UserServiceHost, server.config.UserServicePort)
	return persistence.NewUserServiceClient(address)
}

func (server *Server) startGrpcServer(connectionHandler *api.ConnectionHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to listen in connection service: ", listener)
		log.Fatalf("failed to listen: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Parsing RSA public key in connection service failed, PK:", server.config.PublicKey)
		log.Fatalf("failed to parse public key: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), publicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))

	inventory.RegisterConnectionServiceServer(grpcServer, connectionHandler)
	if err := grpcServer.Serve(listener); err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to serve gRPC in connection service: ", listener)
		log.Fatalf("failed to serve: %s", err)
	}
}
