package main

import (
	"fmt"

	startup "github.com/XWS-BSEP-TIM5-2022/xws-bsep/tree/feat/user-servicesep/microservices/user_service/startup"
	cfg "github.com/XWS-BSEP-TIM5-2022/xws-bsep/tree/feat/user-servicesep/microservices/user_service/startup/config"
)

func main() {
	fmt.Println("Hello world from user")

	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
