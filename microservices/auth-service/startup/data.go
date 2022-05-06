package startup

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
)

var auths = []*domain.Authentication{
	{
		Id:       "1",
		Username: "Ranko",
		Password: "pass2",
		Role:     "ADMIN",
	},
	{
		Id:       "2",
		Username: "Marko",
		Password: "pass2",
		Role:     "REGISTERED_USER",
	},
}
