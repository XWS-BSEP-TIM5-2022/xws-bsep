package startup

import (
	"fmt"
	"log"
	"net"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	"github.com/dgrijalva/jwt-go"

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
}

func NewServer(config *config.Config) *Server {
	CustomLogger := api.NewCustomLogger()
	return &Server{
		config:       config,
		CustomLogger: CustomLogger,
	}
}

const (
	QueueGroup = "connection_service"
)

func (server *Server) Start() {
	neo4jClient := server.initNeo4J()
	server.CustomLogger.SuccessLogger.Info("Neo4J initialization for connection service successful")

	connectionStore := server.initConnectionStore(neo4jClient)
	connectionService := server.initConnectionService(connectionStore)
	connectionHandler := server.initConnectionHandler(connectionService)

	server.CustomLogger.SuccessLogger.Info("Starting gRPC server for connection service")
	server.startGrpcServer(connectionHandler)
}

func (server *Server) initNeo4J() *neo4j.Driver {

	//uri := "bolt:\\" + server.config.ConnectionDBHost + ":" + server.config.ConnectionDBPort
	dbUri := "bolt://localhost:7687"

	client, err := persistence.GetClient(dbUri, server.config.Neo4jUsername, server.config.Neo4jPassword)
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

func (server *Server) initConnectionHandler(service *application.ConnectionService) *api.ConnectionHandler {
	return api.NewConnectionHandler(service)
}

func (server *Server) startGrpcServer(connectionHandler *api.ConnectionHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Starting gRPC server for connection service failed")
		log.Fatalf("failed to listen: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Parsing RSA public key for connection service failed")
		log.Fatalf("failed to parse public key: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), publicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))

	inventory.RegisterConnectionServiceServer(grpcServer, connectionHandler)
	if err := grpcServer.Serve(listener); err != nil {
		server.CustomLogger.ErrorLogger.Error("Serving gRPC server for connection service failed")
		log.Fatalf("failed to serve: %s", err)
	}
}
