package main

import (
	"fmt"

	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth-service/startup"

	cfg "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/auth_service/startup/config"
)

func main() {
	fmt.Println("Hello world from auth")

	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()

}
