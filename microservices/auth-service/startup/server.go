package startup

import (
	"fmt"
	"log"
	"net"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth_service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth_service/startup/config"
	auth_service_proto "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
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
	fmt.Println("Postgres cli")
	fmt.Println(postgresClient)
	authStore := server.initAuthStore(postgresClient)
	authService := server.initAuthService(authStore)

	authHandler := server.initAuthHandler(authService)

	server.startGrpcServer(authHandler)
}

func (server *Server) initPostgresClient() *gorm.DB {
	client, err := persistence.GetClient(
		server.config.AuthDBHost, server.config.AuthDBAuth,
		server.config.AuthDBPass, server.config.AuthDBName,
		server.config.AuthDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initAuthStore(client *gorm.DB) domain.AuthStore {
	store, err := persistence.NewAuthPostgresStore(client)
	if err != nil {
		log.Fatal(err)
	}
	store.DeleteAll()
	// for _, Auth := range auths {
	// 	res, err := store.Insert(Auth)
	// 	fmt.Println(res)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	return store
}

func (server *Server) initAuthService(store domain.AuthStore) *application.AuthService {
	return application.NewAuthService(store)
}

func (server *Server) initAuthHandler(service *application.AuthService) *api.AuthHandler {
	return api.NewAuthHandler(service)
}

func (server *Server) startGrpcServer(authHandler *api.AuthHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	auth_service_proto.RegisterAuthServiceServer(grpcServer, authHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
