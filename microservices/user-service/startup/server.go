package startup

import (
	"fmt"
	"log"
	"net"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/infrastructure/persistence"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"

	interceptor "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/interceptor"
	inventory "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/infrastructure/api"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/startup/config"
	"google.golang.org/grpc"
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
	QueueGroup = "user_service"
)

func (server *Server) Start() {
	//postgresClient := server.initPostgresClient()
	//userStore := server.initUserStore(postgresClient)

	mongoClient := server.initMongoClient()
	userStore := server.initUserStore(mongoClient)

	userService := server.initUserService(userStore)

	userHandler := server.initUserHandler(userService)

	server.startGrpcServer(userHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.UserDBHost, server.config.UserDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initUserStore(client *mongo.Client) domain.UserStore {
	store := persistence.NewUserMongoDBStore(client)
	store.DeleteAll()
	for _, user := range users {
		_, err := store.Insert(user)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

//func (server *Server) initPostgresClient() *gorm.DB {
//	client, err := persistence.GetClient(
//		server.config.UserDBHost, server.config.UserDBUser,
//		server.config.UserDBPass, server.config.UserDBName,
//		server.config.UserDBPort)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return client
//}
//
//func (server *Server) initUserStore(client *gorm.DB) domain.UserStore {
//	store, err := persistence.NewUserPostgresStore(client)
//	if err != nil {
//		log.Fatal(err)
//	}
//	store.DeleteAll()
//	for _, User := range users {
//		res, err := store.Insert(User)
//		fmt.Println(res)
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//	return store
//}

func (server *Server) initUserService(store domain.UserStore) *application.UserService {
	return application.NewUserService(store)
}

func (server *Server) initUserHandler(service *application.UserService) *api.UserHandler {
	return api.NewUserHandler(service)
}

func (server *Server) startGrpcServer(userHandler *api.UserHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// grpcServer := grpc.NewServer()
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		log.Fatalf("failed to parse public key: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessibleRoles(), publicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	inventory.RegisterUserServiceServer(grpcServer, userHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
