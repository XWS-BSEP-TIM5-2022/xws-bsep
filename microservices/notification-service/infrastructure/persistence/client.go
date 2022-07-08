package persistence

import (
	"context"
	"fmt"
	auth "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	connection "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	message "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/message_service"
	post "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func GetClient(host, port string) (*mongo.Client, error) {

	uri := fmt.Sprintf("mongodb://%s:%s/", host, port)
	options := options.Client().ApplyURI(uri)
	return mongo.Connect(context.TODO(), options)
}

func NewUserServiceClient(address string) user.UserServiceClient {
	con, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC user to notification service: %v", err)
	}
	return user.NewUserServiceClient(con)
}

func NewConnectionServiceClient(address string) connection.ConnectionServiceClient {
	con, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to notification service: %v", err)
	}
	return connection.NewConnectionServiceClient(con)
}

func NewMessageServiceClient(address string) message.MessageServiceClient {
	con, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC message to notification service: %v", err)
	}
	return message.NewMessageServiceClient(con)
}

func NewPostServiceClient(address string) post.PostServiceClient {
	con, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC post to notification service: %v", err)
	}
	return post.NewPostServiceClient(con)
}

func NewAuthServiceClient(address string) auth.AuthServiceClient {
	con, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Auth service: %v", err)
	}
	return auth.NewAuthServiceClient(con)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
