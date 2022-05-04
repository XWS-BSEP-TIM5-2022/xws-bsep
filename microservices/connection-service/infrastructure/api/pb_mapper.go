package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/connection_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/connection_service/domain"
)

func mapUserConn(userConn *domain.UserConn) *pb.User {
	userConnPb := &pb.User{
		UserID:   userConn.UserID,
		IsPublic: userConn.IsPublic,
	}

	return userConnPb
}
