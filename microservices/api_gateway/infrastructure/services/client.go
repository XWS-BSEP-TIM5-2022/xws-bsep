package services

import (
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func NewUserClient(address string) user.UserServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return user.NewUserServiceClient(conn)
}

//func NewPostClient(address string) post.PostServiceClient {
//	conn, err := getConnection(address) // dobavljanje konekcije
//	if err != nil {
//		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
//	}
//	return post.NewPostServiceClient(conn) // kreiran novi gRPC klijent (u odnosu na dobavljenu konekciju)
//}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
