package startup

import "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth_service/startup/config"

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}
