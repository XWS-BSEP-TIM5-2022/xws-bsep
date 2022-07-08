package services

import (
	"log"

	connection "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	message "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/message_service"
	notification "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	post "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/post_service"

	auth "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
	joboffer "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/job_offer_service"
	user "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(address string) user.UserServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to User service: %v", err)
	}
	return user.NewUserServiceClient(conn)
}

func NewAuthClient(address string) auth.AuthServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Auth service: %v", err)
	}
	return auth.NewAuthServiceClient(conn)
}

func NewPostClient(address string) post.PostServiceClient {
	conn, err := getConnection(address) // dobavljanje konekcije
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return post.NewPostServiceClient(conn) // kreiran novi gRPC klijent (u odnosu na dobavljenu konekciju)
}

func NewConnectionClient(address string) connection.ConnectionServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Connection service: %v", err)
	}
	return connection.NewConnectionServiceClient(conn)
}

func NewNotificationClient(address string) notification.NotificationServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to notification service: %v", err)
	}
	return notification.NewNotificationServiceClient(conn)
}

func NewJobOfferClient(address string) joboffer.JobOfferServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Connection service: %v", err)
	}
	return joboffer.NewJobOfferServiceClient(conn)
}

func NewMessageClient(address string) message.MessageServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Message service: %v", err)
	}
	return message.NewMessageServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}
