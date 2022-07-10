package startup

import (
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/tracer"
	otgo "github.com/opentracing/opentracing-go"
	"io"
	"log"
	"net"

	"github.com/dgrijalva/jwt-go"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	inventory "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/job_offer_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/job_offer_service/startup/config"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"google.golang.org/grpc"
)

type Server struct {
	config       *config.Config
	CustomLogger *api.CustomLogger
	tracer       otgo.Tracer
	closer       io.Closer
}

const name = "job-offer-service"

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
	server.CustomLogger.SuccessLogger.Info("Neo4J initialization for job offer service successful, PORT: ", server.config.Port, ", HOST: ", server.config.Host)

	connectionStore := server.initJobOfferStore(neo4jClient)
	connectionService := server.initJobOfferService(connectionStore)
	connectionHandler := server.initJobOfferHandler(connectionService)

	server.CustomLogger.SuccessLogger.Info("Starting gRPC server for job offer service")
	server.startGrpcServer(connectionHandler)
}

func (server *Server) initNeo4J() *neo4j.Driver {

	//uri := "bolt:\\" + server.config.ConnectionDBHost + ":" + server.config.ConnectionDBPort
	// dbUri := "bolt://localhost:7687"
	dbUri := "bolt://" + server.config.ConnectionDBHost + ":" + server.config.ConnectionDBPort

	client, err := persistence.GetClient(dbUri, server.config.ConnectionDBUser, server.config.ConnectionDBPass)
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Neo4J initialization for job offer service failed")
		log.Fatal(err)
	}
	return client
}

func (server *Server) initJobOfferStore(client *neo4j.Driver) domain.JobOfferStore {
	store := persistence.NewJobOfferDBStore(client)
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

func (server *Server) initJobOfferService(store domain.JobOfferStore) *application.JobOfferService {
	return application.NewJobOfferService(store)
}

func (server *Server) initJobOfferHandler(service *application.JobOfferService) *api.JobOfferHandler {
	return api.NewJobOfferHandler(service)
}

func (server *Server) startGrpcServer(jobOfferHandler *api.JobOfferHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to listen in connection service: ", listener)
		log.Fatalf("failed to listen: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		server.CustomLogger.ErrorLogger.Error("Parsing RSA public key in connection service failed")
		log.Fatalf("failed to parse public key: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessiblePermissions(), publicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))

	inventory.RegisterJobOfferServiceServer(grpcServer, jobOfferHandler)
	if err := grpcServer.Serve(listener); err != nil {
		server.CustomLogger.ErrorLogger.Error("Failed to serve gRPC in connection service: ", listener)
		log.Fatalf("failed to serve: %s", err)
	}
}
