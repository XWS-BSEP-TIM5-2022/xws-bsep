package startup

import (
	"time"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/domain"
)

var auths = []*domain.Authentication{
	{
		Id:       "1",
		Name:     "Ranko",
		Password: "pass2",
		Role:     "ADMIN",
		Date:     time.Date(2022, 10, 11, 8, 4, 0, 0, time.UTC),
	},
	{
		Id:       "2",
		Name:     "Marko",
		Password: "pass2",
		Role:     "REGISTERED_USER",
		Date:     time.Now(),
	},
}
