package startup

import (
	"fmt"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/application"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/domain"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/infrastructure/persistence"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/post_service/startup/config"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
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

	fmt.Println(postService.GetAll())

	//postHandler := server.initProductHandler(postService)
	//
	//server.startGrpcServer(postHandler)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.PostDBHost, server.config.PostDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initPostStore(client *mongo.Client) domain.PostStore {
	store := persistence.NewPostMongoDBStore(client)
	store.DeleteAll()
	for _, post := range posts {
		err := store.Insert(post)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) initPostService(store domain.PostStore) *application.PostService {
	return application.NewPostService(store)
}

//func (server *Server) initPostHandler(service *application.PostService) *api.PostHandler {
//	return api.NewPostHandler(service)
//}
//
//func (server *Server) startGrpcServer(productHandler *api.ProductHandler) {
//	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//	grpcServer := grpc.NewServer()
//	catalogue.RegisterCatalogueServiceServer(grpcServer, productHandler)
//	if err := grpcServer.Serve(listener); err != nil {
//		log.Fatalf("failed to serve: %s", err)
//	}
//}
