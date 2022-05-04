package startup

import "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"

var auths = []*domain.Authentication{
	{
		Id:       "1",
		Name:     "Ranko",
		Password: "pass2",
	},
	{
		Id:       "2",
		Name:     "Marko",
		Password: "pass2",
	},
}
