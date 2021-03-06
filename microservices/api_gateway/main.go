package main

import (
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/startup"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway/startup/config"
)

func main() {
	config := config.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
