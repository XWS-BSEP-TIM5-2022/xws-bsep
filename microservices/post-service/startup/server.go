package startup

import (
	"fmt"
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
	"log"
	"net"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	mongoClient := server.initMongoClient()
	postStore := server.initPostStore(mongoClient)
	postService := server.initPostService(postStore)
	postHandler := server.initPostHandler(postService)

	server.startGrpcServer(postHandler)
}

func (server *Server) initMongoClient() *mongo.Client { // inicijalizacija mongo klijenta
	client, err := persistence.GetClient(server.config.PostDBHost, server.config.PostDBPort)
	if err != nil {
		log.Fatal(err)
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

func (server *Server) initPostService(store domain.PostStore) *application.PostService {
	return application.NewPostService(store)
}

func (server *Server) initPostHandler(service *application.PostService) *api.PostHandler { // inicjializacija post handler-a, odnosno post servisa
	return api.NewPostHandler(service)
}

func (server *Server) startGrpcServer(postHandler *api.PostHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(server.config.PublicKey))
	if err != nil {
		log.Fatalf("failed to parse public key: %v", err)
	}

	interceptor := interceptor.NewAuthInterceptor(config.AccessibleRoles(), publicKey)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()))
	post.RegisterPostServiceServer(grpcServer, postHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
