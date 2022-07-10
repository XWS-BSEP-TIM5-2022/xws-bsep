package startup

import (
	"fmt"
	event "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/event_service"
	"github.com/sirupsen/logrus"
	"log"
	"net"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/event_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/event_service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/event_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/event_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/event_service/startup/config"
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
	server.CustomLogger.SuccessLogger.Info("MongoDB initialization for event service successful, PORT: ", server.config.EventDBPort, ", HOST: ", server.config.EventDBHost)

	eventStore := server.initEventStore(mongoClient)
	eventService := server.initEventService(eventStore)
	eventHandler := server.initEventHandler(eventService)

	server.CustomLogger.SuccessLogger.Info("Starting gRPC server for event service")
	server.startGrpcServer(eventHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.EventDBHost, server.config.EventDBPort)
	if err != nil {
		server.CustomLogger.ErrorLogger.WithFields(logrus.Fields{
			"event_db_host": server.config.EventDBHost,
			"event_db_port": server.config.EventDBPort,
		}).Error("MongoDB initialization for event service failed")
	}
	return client
}

func (server *Server) initEventStore(client *mongo.Client) domain.EventStore {
	store := persistence.NewEventMongoDBStore(client)
	//store.DeleteAll()
	//for _, event := range events {
	//	_, err := store.Insert(event)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	return store
}

func (server *Server) initEventService(store domain.EventStore) *application.EventService {
	return application.NewEventService(store)
}

func (server *Server) initEventHandler(service *application.EventService) *api.EventHandler {
	return api.NewEventHandler(service)
}

func (server *Server) startGrpcServer(eventHandler *api.EventHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to listen in event service: ", listener)
		log.Fatalf("failed to listen: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Parsing RSA public key in event service failed")
		log.Fatalf("failed to parse public key: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), publicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))

	event.RegisterEventServiceServer(grpcServer, eventHandler)
	if err := grpcServer.Serve(listener); err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to serve gRPC in event service: ", listener)
		log.Fatalf("failed to serve: %s", err)
	}
}
