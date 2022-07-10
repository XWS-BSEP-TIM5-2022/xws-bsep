package persistence

import (
	"context"  
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"fmt"
	notification "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func GetClient(uri, username, password string) (*neo4j.Driver, error) {

	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &driver, nil
}

func GetMongoClient(host, port string) (*mongo.Client, error) {

	uri := fmt.Sprintf("mongodb://%s:%s/", host, port)
	options := options.Client().ApplyURI(uri)
	return mongo.Connect(context.TODO(), options)
}
func NewNotificationServiceClient(address string) notification.NotificationServiceClient {
	con, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Notification service: %v", err)
	}
	return notification.NewNotificationServiceClient(con)
}

func NewUserServiceClient(address string) user.UserServiceClient {
	con, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to User service: %v", err)
	}
	return user.NewUserServiceClient(con)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

}
