package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/user_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/user_service/domain"
)

func mapUser(user *domain.User) *pb.User {
	userPb := &pb.User{
		Id:   user.Id,
		Name: user.Name,
	}
	return userPb
}
