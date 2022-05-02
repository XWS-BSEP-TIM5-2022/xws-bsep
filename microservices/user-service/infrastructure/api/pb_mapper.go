package api

import (
	pb "github.com/sanjadrinic/test_repo/microservices/common/proto/user_service"
	"github.com/sanjadrinic/test_repo/microservices/user_service/domain"
)

func mapUser(order *domain.User) *pb.User {
	userPb := &pb.User{
		Id:   order.Id,
		Name: order.Name,
	}
	return userPb
}
