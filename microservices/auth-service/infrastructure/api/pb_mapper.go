package api

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/auth_service"
)

func mapNewAuth(request *pb.RegisterRequest) *domain.Authentication {
	auth := &domain.Authentication{
		Id:       "", // TODO SD: !!!!!!!!! ovo treba ispraviti nakon uspesno dodatog user-a u user servisu
		Username: request.Username,
		Password: request.Password,
		Roles:    &[]domain.Role{},
	}
	return auth
}
