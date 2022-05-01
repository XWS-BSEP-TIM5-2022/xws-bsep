package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/tree/feat/user-servicesep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/tree/feat/user-servicesep/microservices/user_service/domain"
)

func mapUser(order *domain.User) *pb.User {
	userPb := &pb.User{
		Id:   order.Id.Hex(),
		Name: order.Name,
	}
	return userPb
}
