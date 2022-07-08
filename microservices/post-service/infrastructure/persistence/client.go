package persistence

import (
	"context"
	"fmt"
	auth "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	notification "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func GetClient(host, port string) (*mongo.Client, error) { // pomocna metoda koja vraca mongo klijenta na osnovu neophodnih podataka (host, port)
	uri := fmt.Sprintf("mongodb://%s:%s/", host, port)
	options := options.Client().ApplyURI(uri)
	return mongo.Connect(context.TODO(), options) // konektuje se i vraca klijenta nazad
}

func NewUserServiceClient(address string) user.UserServiceClient {
	con, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to User service: %v", err)
	}
	return user.NewUserServiceClient(con)
}

func NewAuthServiceClient(address string) auth.AuthServiceClient {
	con, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Auth service: %v", err)
	}
	return auth.NewAuthServiceClient(con)
}

func NewNotificationServiceClient(address string) notification.NotificationServiceClient {
	con, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Notification service: %v", err)
	}
	return notification.NewNotificationServiceClient(con)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
