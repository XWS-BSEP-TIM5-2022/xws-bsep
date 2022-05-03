package api

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
)

func mapCreateAuth(authPb *pb.Auth) *domain.Authentication {
	auth := &domain.Authentication{
		Id:       authPb.Id,
		Name:     authPb.Name,
		Password: authPb.Password,
	}

	return auth
}
